import sqlite3

db = sqlite3.connect('buzz.sqlite', check_same_thread=False)  #: dangerous: "As long as only as a single thread is writing through the connection in a given time, this is safe to use."
