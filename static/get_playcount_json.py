import sqlite3
import json

# Connect to the SQLite database
connection = sqlite3.connect('flashpoint.sqlite')

# Create a cursor object using the cursor() method
cursor = connection.cursor()

# SQL query to select id and playCounter from the game table where playCounter is larger than 0
sql_query = "SELECT id, playCounter FROM game WHERE playCounter > 0"

# Execute the query
cursor.execute(sql_query)

# Fetch all rows
rows = cursor.fetchall()

# Close the connection
connection.close()

# Convert rows to the required JSON format
played_games = {row[0]: row[1] for row in rows}

# Write to a JSON file
with open('games_played.json', 'w') as json_file:
    json.dump(played_games, json_file, indent=4)

print("Data successfully written to games_played.json")