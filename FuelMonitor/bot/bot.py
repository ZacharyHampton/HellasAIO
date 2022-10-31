import nextcord
from nextcord.abc import GuildChannel
from nextcord import Interaction, SlashOption, CategoryChannel
from nextcord.ext import application_checks
from nextcord.ext import commands
import sqlite3
from fuel.utils import DontRunKeywords, CurrentRunningKeywords
from db import db
import time


class DiscordHandler(commands.Cog):
    def __init__(self, bot: commands.Bot):
        self.bot = bot

    @commands.Cog.listener()
    async def on_ready(self):
        print("{} is ready!".format(self.bot.user.name))

    @commands.Cog.listener()
    async def on_command_error(self, ctx: Interaction, error):
        if isinstance(error, commands.MissingPermissions):
            await ctx.response.send_message("You don't have permission to use this command.")
            return

        raise error

    @nextcord.slash_command(description="Main bot controlling commands")
    @application_checks.has_permissions(manage_messages=True)
    async def fuel(self, interaction: Interaction):
        """rows = db.cursor().execute('SELECT discord_id FROM permitted_users').fetchall()

        if interaction.message.author.id not in [row[0] for row in rows]:
            await interaction.response.send_message("You are not permitted to use this command.")
            return"""

    @fuel.subcommand(description="Add keyword to fuel monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def add_keyword(
            self,
            interaction: Interaction,
            keyword: str = SlashOption(name="keyword", description="Product Keyword"),
    ):
        cur = db.cursor()
        if keyword in DontRunKeywords:
            DontRunKeywords.remove(keyword)

        cur.execute("INSERT INTO keywords VALUES (?)", (keyword,))
        db.commit()

        await interaction.response.send_message("Keyword added.")

    @fuel.subcommand(description="Add keyword to fuel monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def delete_keyword(
            self,
            interaction: Interaction,
            keyword: str = SlashOption(name="keyword", description="Product Keyword"),
    ):
        cur = db.cursor()
        if keyword not in DontRunKeywords:
            DontRunKeywords.append(keyword)

        if keyword in CurrentRunningKeywords:
            CurrentRunningKeywords.remove(keyword)

        cur.execute("DELETE FROM keywords WHERE keyword = ?", (keyword,))
        #: wait for flow to cycle out of above kw
        time.sleep(5)

        cur.execute("DELETE FROM keyword_data WHERE keyword = ?", (keyword,))
        db.commit()

        await interaction.response.send_message("Keyword removed.")

    @fuel.subcommand(description="Add MSKU to fuel monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def add_msku(
            self,
            interaction: Interaction,
            msku: str = SlashOption(name="msku", description="Product's Manufacture SKU"),
    ):
        cur = db.cursor()
        cur.execute("INSERT INTO mskus VALUES (?)", (msku,))
        db.commit()

        await interaction.response.send_message("MSKU added.")

    @fuel.subcommand(description="Remove MSKU to fuel monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def delete_msku(
            self,
            interaction: Interaction,
            msku: str = SlashOption(name="msku", description="Product's Manufacture SKU"),
    ):
        cur = db.cursor()
        cur.execute("DELETE FROM mskus WHERE msku = ?", (msku,))
        db.commit()

        await interaction.response.send_message("MSKU removed.")

        #: wait for flow to cycle out of above msku
        time.sleep(5)
        cur.execute("DELETE FROM msku_data WHERE sku = ?", (msku + '%',))
        db.commit()


def setup(bot):
    bot.add_cog(DiscordHandler(bot))
