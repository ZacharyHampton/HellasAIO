import sys
import nextcord
from nextcord.ext import commands
import os
from db import db
from buzz.start import Start
import threading
import sentry_sdk
from dotenv import load_dotenv

load_dotenv()


def main():
    sentry_sdk.init(dsn=os.getenv('SENTRY_DSN'))

    intents = nextcord.Intents.all()
    extensions = ['bot.bot']
    bot = commands.Bot(
        command_prefix='!',
        intents=intents,
        rollout_guild_ids=[int(os.getenv('GUILD_ID'))],
    )

    for ext in extensions:
        bot.load_extension(ext)

    threading.Thread(target=Start).start()
    bot.run(os.getenv('DISCORD_TOKEN'), reconnect=True)

    while True:
        q = input()
        if q.lower() in ['q', 'quit']:
            db.close()
            sys.exit(0)


if __name__ == '__main__':
    main()
