import requests
import re
import bypass.get_session as get_session


def get_buildnum():
    session = get_session.get_session()
    text = session.get("https://discord.com/login").text 
    script_url = 'https://discord.com/assets/' + re.compile(r'\d+\.\w+\.js|sentry\.\w+\.js').findall(text)[-1]
    text = session.get(script_url).text
    index = text.find("buildNumber") + 26
    build_num = int(text[index:index + 6])
    return build_num