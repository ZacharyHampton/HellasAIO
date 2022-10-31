import os
import requests
import platform
import sys
import hashlib
import subprocess


def getChecksum():
    sha256_hash = hashlib.sha256()

    if not os.path.exists('bin/bot.exe'):
        if not os.path.exists('bin'):
            os.mkdir('bin')

        return ""

    with open("bin/bot.exe", "rb") as f:
        for byte_block in iter(lambda: f.read(4096), b""):
            sha256_hash.update(byte_block)

    return sha256_hash.hexdigest()


def main():
    print("Welcome to HellasAIO!")
    if platform.system() != "Windows":
        print('We only support windows currently.')
        sys.exit(0)

    print("Checking for updates...")
    response = requests.get('https://api.hellasaio.com/api/latest')
    if response.status_code != 200 or not response.json()['success']:
        print('Failed to check for updates.')
        sys.exit(0)

    if response.json()['downloads']['windows']['checksum'] != getChecksum():
        print('Checksum of new file:', response.json()['downloads']['windows']['checksum'])
        print('Please compare the checksum here with the one in the #changelogs channel to prevent any supply chain attacks. Downloading this without properly comparing could result in your computer being compromised.')
        c = input('Update found. Do you want to update? (y/n) ')
        if c.lower() == 'n':
            print('Update aborted.')
            sys.exit(0)
        elif c.lower() != 'y':
            print('Invalid input.')
            sys.exit(0)

        print('Downloading new version...')
        response = requests.get(response.json()['downloads']['windows']['url'])
        if response.status_code != 200:
            print('Failed to download new version.')
            sys.exit(0)
        with open("bin/bot.exe", "wb") as f:
            f.write(response.content)
        print('Downloaded new version.')
    else:
        print('Bot is up to date!')

    print('Starting bot...')
    os.system('cls')
    os.system(r'cd bin && bot.exe')
    input('Press enter to exit...')
    sys.exit(0)


if __name__ == '__main__':
    main()
