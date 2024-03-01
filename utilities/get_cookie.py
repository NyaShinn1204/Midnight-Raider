def get_cookie(session):
    response = session.get("https://discord.com/api/v9/experiments").cookies
    return f'__dcfduid={response["__dcfduid"]}; __sdcfduid={response["__sdcfduid"]}; locale=ja-JP'