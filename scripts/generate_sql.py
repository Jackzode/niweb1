import mysql.connector
from mysql.connector import Error
import random
import datetime

def create_connection(host_name, port, user_name, user_password, db_name):
    connection = None
    try:
        connection = mysql.connector.connect(
            host=host_name,
            user=user_name,
            port = port,
            passwd=user_password,
            database=db_name,
        )
        print("Connection to MySQL DB successful")
    except Error as e:
        print(f"The error '{e}' occurred")
    return connection

def insert_question(connection, question):
    cursor = connection.cursor()
    try:
        cursor.execute(question)
        connection.commit()
        print("Query executed successfully")
    except Error as e:
        print(f"The error '{e}' occurred")

def generate_insert_query(id):
    created_at = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    user_id = random.randint(1, 100)
    last_edit_user_id = random.randint(1, 100)
    title = f'Title {id}'
    original_text = f'<p>Original text for question {id}</p>'
    parsed_text = f'<p>Parsed text for question {id}</p>'
    pin = random.choice([0, 1])
    show = 1
    status = 1
    view_count = random.randint(0, 1000)
    unique_view_count = random.randint(0, 1000)
    vote_count = random.randint(0, 1000)
    answer_count = random.randint(0, 1000)
    collection_count = random.randint(0, 1000)
    follow_count = random.randint(0, 1000)
    post_update_time = created_at
    revision_id = random.randint(0, 100)
    accepted_answer_id = random.randint(0, 1000)
    last_answer_id = random.randint(0, 1000)
    copyright = 1
    allow_reprint = 1
    allow_comment = 1
    feeds = 0

    query = f"""
    INSERT INTO `question`
    (`id`, `created_at`, `updated_at`, `user_id`, `invite_user_id`, `last_edit_user_id`, `title`, `original_text`, `parsed_text`, `pin`, `show`, `status`, `view_count`, `unique_view_count`, `vote_count`, `answer_count`, `collection_count`, `follow_count`, `accepted_answer_id`, `last_answer_id`, `post_update_time`, `revision_id`, `copyright`, `allow_reprint`, `allow_comment`, `feeds`)
    VALUES
    ({id}, '{created_at}', NULL, {user_id}, NULL, {last_edit_user_id}, '{title}', '{original_text}', '{parsed_text}', {pin}, {show}, {status}, {view_count}, {unique_view_count}, {vote_count}, {answer_count}, {collection_count}, {follow_count}, {accepted_answer_id}, {last_answer_id}, '{post_update_time}', {revision_id}, {copyright}, {allow_reprint}, {allow_comment}, {feeds});
    """
    return query

def main():
    connection = create_connection("localhost","3306", "root", "root", "lawyer")

    for i in range(1, 1001):
        query = generate_insert_query(10010000001000000 + i)
        insert_question(connection, query)

if __name__ == "__main__":
    main()
