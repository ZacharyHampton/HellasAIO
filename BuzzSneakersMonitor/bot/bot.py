import nextcord
from nextcord.abc import GuildChannel
from nextcord import Interaction, SlashOption, CategoryChannel
from nextcord.ext import application_checks
from nextcord.ext import commands
from buzz.threads import RunningThreads
import sqlite3
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
    async def buzz(self, interaction: Interaction):
        """rows = db.cursor().execute('SELECT discord_id FROM permitted_users').fetchall()

        if interaction.message.author.id not in [row[0] for row in rows]:
            await interaction.response.send_message("You are not permitted to use this command.")
            return"""

    @buzz.subcommand(description="Add Product ID to the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def add_product_id(
            self,
            interaction: Interaction,
            pid: str = SlashOption(name="pid", description="Product's ID"),
    ):
        if pid in RunningThreads.keys():
            await interaction.response.send_message("This product is already being monitored.")
            return

        cur = db.cursor()
        cur.execute("INSERT INTO products VALUES (?)", (pid,))
        db.commit()

        await interaction.response.send_message("PID added.")

    @buzz.subcommand(description="Remove Product ID from the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def delete_product_id(
            self,
            interaction: Interaction,
            pid: str = SlashOption(name="pid", description="Product's ID"),
    ):
        if pid in RunningThreads.keys():
            RunningThreads[pid].Stop()

        cur = db.cursor()
        cur.execute("DELETE FROM products WHERE id = ?", (pid,))
        db.commit()

        await interaction.response.send_message("PID removed.")

        #: wait for flow to cycle out of above msku
        time.sleep(5)
        cur.execute("DELETE FROM product_data WHERE id = ?", (pid,))
        db.commit()

    @buzz.subcommand(description="Add MSKU to the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def add_msku(
            self,
            interaction: Interaction,
            msku: str = SlashOption(name="msku", description="Manufacture SKU"),
    ):
        if msku in RunningThreads.keys():
            await interaction.response.send_message("This product is already being monitored.")
            return

        cur = db.cursor()
        cur.execute("INSERT INTO mskus VALUES (?)", (msku,))
        db.commit()

        await interaction.response.send_message("MSKU added.")

    @buzz.subcommand(description="Remove MSKU from the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def delete_msku(
            self,
            interaction: Interaction,
            msku: str = SlashOption(name="msku", description="Manufacture SKU"),
    ):
        if msku in RunningThreads.keys():
            RunningThreads[msku].Stop()

        cur = db.cursor()
        cur.execute("DELETE FROM mskus WHERE msku = ?", (msku,))
        db.commit()

        await interaction.response.send_message("MSKU removed.")

    @buzz.subcommand(description="Add secret to the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def add_backend(
            self,
            interaction: Interaction,
            backend: str = SlashOption(name="backend", description="secret"),
    ):
        if backend in RunningThreads.keys():
            await interaction.response.send_message("This product is already being monitored.")
            return

        cur = db.cursor()
        cur.execute("INSERT INTO msku_partials_image VALUES (?)", (backend,))
        db.commit()

        await interaction.response.send_message("secret added.")

    @buzz.subcommand(description="Remove secret from the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def delete_backend(
            self,
            interaction: Interaction,
            backend: str = SlashOption(name="backend", description="secret"),
    ):
        if backend in RunningThreads.keys():
            RunningThreads[backend].Stop()

        cur = db.cursor()
        cur.execute("DELETE FROM msku_partials_image WHERE msku_partial = ?", (backend,))
        db.commit()

        await interaction.response.send_message("secret removed.")

    @buzz.subcommand(description="Add keyword to the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def add_keyword(
            self,
            interaction: Interaction,
            keyword: str = SlashOption(name="keyword", description="keyword"),
    ):
        if keyword in RunningThreads.keys():
            await interaction.response.send_message("This product is already being monitored.")
            return

        cur = db.cursor()
        cur.execute("INSERT INTO keywords VALUES (?)", (keyword,))
        db.commit()

        await interaction.response.send_message("keyword added.")

    @buzz.subcommand(description="Remove keyword from the buzz monitor.")
    @application_checks.has_permissions(manage_messages=True)
    async def delete_keyword(
            self,
            interaction: Interaction,
            keyword: str = SlashOption(name="keyword", description="keyword"),
    ):
        if keyword in RunningThreads.keys():
            RunningThreads[keyword].Stop()

        cur = db.cursor()
        cur.execute("DELETE FROM keywords WHERE keyword = ?", (keyword,))
        db.commit()

        await interaction.response.send_message("keyword removed.")


def setup(bot):
    bot.add_cog(DiscordHandler(bot))
