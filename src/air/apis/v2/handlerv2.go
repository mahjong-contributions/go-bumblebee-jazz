package v2

import (
	"air/aqi"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	RedisServer         string
	RedisServerPassword string
	RedisClient         *redis.Client
)

type ResponseAirQuality struct {
	ServerVersion string         `json:"server_version"`
	Air           aqi.AirQuality `json:"air_quality"`
}

func headers(c *gin.Context) {
	ver := os.Getenv("AIR_VERSION")
	if ver == "" || !strings.HasPrefix(ver, "v2") {
		ver = "v2"
	}
	c.Header("air_server", "air")
	c.Header("air_version", ver)
}

func init() {
	RedisServer = os.Getenv("REDIS_SERVER_ADDRESS")
	RedisServerPassword = os.Getenv("REDIS_SERVER_PASSWORD")
	if RedisServer == "" {
		log.Error("Initial redis server address was failed.")
	} else {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     RedisServer,
			Password: RedisServerPassword,
		})
		pong, err := RedisClient.Ping(context.TODO()).Result()
		if err != nil {
			log.Error(pong, err)
		}
	}
}

func CachingAirQuality(ctx context.Context, city string, quality aqi.AirQuality) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "CachingAirQuality")
	defer span.Finish()

	buf, _ := json.Marshal(&quality)
	if RedisClient == nil {
		return errors.New("redis client was nil")
	}
	pipeline := RedisClient.Pipeline()
	pipeline.Set(ctx, "air_quality_cache-"+city, "expired-600s", 600*time.Second)
	pipeline.HSet(ctx, "air_quality_cache", city, buf)
	_, err := pipeline.Exec(ctx)

	return err
}

func CachedAirQuality(ctx context.Context, city string) (aqi.AirQuality, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "CachedAirQuality")
	defer span.Finish()

	var air = aqi.AirQuality{}

	if RedisClient == nil {
		return air, errors.New("redis client was nil")
	}
	r, err := RedisClient.Get(ctx, "air_quality_cache-"+city).Result()
	log.Infof("air_quality_cache-%s = %s", city, r)
	if err == nil {
		buf, err := RedisClient.HGet(ctx, "air_quality_cache", city).Bytes()
		if err != nil {
			return air, err
		}
		err = json.Unmarshal(buf, &air)
	} else {
		err = errors.New("Cached Air Quality of " + city + "has been expired.")
	}

	return air, err

}

func AirOfGeo(ctx context.Context, c *gin.Context) {
	///air/geo/:lat/:lng ->//feed/geo::lat;:lng/?token=:token
	//Auckland: -36.916839599609375, 174.70875549316406
	headers(c)
	span, _ := opentracing.StartSpanFromContext(ctx, "http-AirOfGeo")
	defer span.Finish()

	lat := c.Param("lat")
	lng := c.Param("lng")
	url := aqi.AQIServer + "/feed/geo:" + lat + ";" + lng + "/?token=" + aqi.AQIServerToken
	if buf, err := aqi.HttpGet(ctx, url); err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)

	} else {
		if air, err := convertAir(buf); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			ver := os.Getenv("AIR_VERSION")
			if ver=="" || !strings.HasPrefix(ver, "v2") {
				c.JSON(http.StatusOK, ResponseAirQuality{
					ServerVersion: "v2",
					Air:           air,
				})
			} else {
				c.JSON(http.StatusOK, ResponseAirQuality{
					ServerVersion: ver,
					Air:           air,
				})
			}
		}
	}
}

func byCity(ctx context.Context, city string) (aqi.AirQuality, error) {

	url := aqi.AQIServer + "/feed/" + city + "/?token=" + aqi.AQIServerToken
	// ---
	buf, err := aqi.HttpGet(ctx, url)
	if err != nil {
		log.Errorf("Fail to call AQIServer service from %s", url)
		return aqi.AirQuality{}, err
	}
	return convertAir(buf)
}

func AirOfCity(ctx context.Context, c *gin.Context) {
	headers(c)
	span, sctx := opentracing.StartSpanFromContext(ctx, "http-AirOfCity")
	defer span.Finish()

	city := c.Param("city")

	air, err := CachedAirQuality(sctx, city)
	if err != nil {
		log.Error(err)
		log.Infof("No cache for %s and looking for fresh value.\n ", city)
		air, err := byCity(sctx, city)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			if err := CachingAirQuality(sctx, city, air); err != nil {
				log.Errorf("Caching air quality data was failed. -> %s, %s\n", air.City, err)
			}
			log.Infof("Air Quality of %s was cached.\n ", city)
			ver := os.Getenv("AIR_VERSION")
			if ver=="" || !strings.HasPrefix(ver, "v2") {
				c.JSON(http.StatusOK, ResponseAirQuality{
					ServerVersion: "v2",
					Air:           air,
				})
			} else {
				c.JSON(http.StatusOK, ResponseAirQuality{
					ServerVersion: ver,
					Air:           air,
				})
			}

		}

	} else {
		log.Infof("Return cached Air Quality of %s.\n ", city)
		c.JSON(http.StatusOK, air)
	}

}

func convertAir(content []byte) (aqi.AirQuality, error) {
	var originAir aqi.OriginAirQuality
	var newAir aqi.AirQuality
	var apiError aqi.ApiError

	err := json.Unmarshal(content, &originAir)
	if err != nil {
		log.Println(err)
		return newAir, err
	}
	if originAir.Status == "error" {
		err = json.Unmarshal(content, &apiError)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Convert data was failed due to <%s>. ", apiError.Data)
		return newAir, err

	}
	newAir = aqi.Copy2AirQuality(originAir)

	return newAir, nil

}

func AirOfIP(ctx context.Context, c *gin.Context) {
	headers(c)
	span, sctx := opentracing.StartSpanFromContext(ctx, "http-AirOfIP")
	defer span.Finish()

	ip := c.Param("ip")
	if city, err := aqi.CityByIP(sctx, ip); err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		if air, err := byCity(sctx, city); err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			ver := os.Getenv("AIR_VERSION")
			if ver=="" || !strings.HasPrefix(ver, "v2") {
				c.JSON(http.StatusOK, ResponseAirQuality{
					ServerVersion: "v2",
					Air:           air,
				})
			} else {
				c.JSON(http.StatusOK, ResponseAirQuality{
					ServerVersion: ver,
					Air:           air,
				})
			}
		}

	}
}

func AQIStandard() string {
	return `
			{ 
				"version": "v2",
				{
					"Standard": "Air Quality Index scale as defined by the US-EPA 2016 standard.DEMO HAHA!!!",
					"Definitions": [
						{
							"AQIServer": "0-50",
							"Level": "Good",
							"Implication": "Air quality is considered satisfactory, and air pollution poses little or no risk",
							"Caution": "None"
						},
						{
							"AQIServer": "51 -100",
							"Level": "Moderate",
							"Implication": "Air quality is acceptable; however, for some pollutants there may be a moderate health concern for a very small number of people who are unusually sensitive to air pollution.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
						},
						{
							"AQIServer": "101-150",
							"Level": "Unhealthy for Sensitive Groups",
							"Implication": "Members of sensitive groups may experience health effects. The general public is not likely to be affected.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should limit prolonged outdoor exertion."
						},
						{
							"AQIServer": "151-200",
							"Level": "Unhealthy",
							"Implication": "Everyone may begin to experience health effects; members of sensitive groups may experience more serious health effects",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid prolonged outdoor exertion; everyone else, especially children, should limit prolonged outdoor exertion"
						},
						{
							"AQIServer": "201-300",
							"Level": "Very Unhealthy",
							"Implication": "Health warnings of emergency conditions. The entire population is more likely to be affected.",
							"Caution": "Active children and adults, and people with respiratory disease, such as asthma, should avoid all outdoor exertion; everyone else, especially children, should limit outdoor exertion."
						},
						{
							"AQIServer": "300+",
							"Level": "Hazardous",
							"Implication": "Health alert: everyone may experience more serious health effects",
							"Caution": "Everyone should avoid all outdoor exertion"
						}
					]
				}
			}
		`
}
