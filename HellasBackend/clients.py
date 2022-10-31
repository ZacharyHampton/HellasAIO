from pymongo import MongoClient as mongoClient
import os
import certifi
import helheim

MongoClient = mongoClient(
    "mongodb+srv://root:{}@mongocluster.com/?retryWrites=true&w=majority".format(
        os.getenv('MONGO_PASS')), tlsCAFile=certifi.where())
helheim.auth(os.getenv('HELHEIM_KEY'))
