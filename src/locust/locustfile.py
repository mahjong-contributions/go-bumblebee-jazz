import random
from locust import HttpUser, task, between

class QuickstartUser(HttpUser):
    wait_time = between(5, 10)

    @task
    def index_page(self):
        self.client.get("/air/v1/version")
        self.client.get("/air/v1/aqi")

    @task(5)
    def task1(self):
        cities = ["beijing", 
                    "chengdu", 
                    "auckland", 
                    "london", 
                    "shanghai", 
                    "tianjing", 
                    "dalian", 
                    "urumqi",
                    "shenzhen",
                    "guangzhou",
                    "wuhan",
                    "xiamen",
                    "chongqing",
                    "guangyuan"]
        id = random.randint(0, len(cities)-1)
        self.client.get("/air/v1/city/"+cities[id], name="/air/v1/city/"+cities[id])

    @task(2)
    def task2(self):
        geo = ["34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751", 
                    "34.68291096793206/109.57214355468751",
                    "34.68291096793206/109.57214355468751",
                    "34.68291096793206/109.57214355468751",
                    "34.68291096793206/109.57214355468751",
                    "34.68291096793206/109.57214355468751"]
        id = random.randint(0, len(geo)-1)
        self.client.get("/air/v1/geo/"+geo[id], name="/air/v1/geo"+geo[id])
