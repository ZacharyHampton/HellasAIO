<template>
  <div class="grid h-screen place-items-center font-sans">
    <LogoImport title="HellasAIO Quicktasking" />

    <div class="">
      <button
          type="button"
          class="flex box-border items-center text-white bg-blurple hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 disabled:bg-slate-50 disabled:text-slate-500 disabled:border-slate-200 font-medium rounded-lg text-sm px-5 py-2.5 mr-2 mb-2"
          @click="redirectToDiscord()"
          ref="discordButton"
      >
        <img class="h-4 pr-2" :src="DiscordLogo"  alt=""/>
        Continue with Discord
      </button>
    </div>
  </div>
</template>


<script>
import DiscordLogo from '@/assets/Discord-Logo-White.svg'
import HellasAIOLogo from '@/assets/logo.png'
import LogoImport from "@/components/LogoImport";
import router from '@/router'

function generateRandomString() {
  let randomString = '';
  const randomNumber = Math.floor(Math.random() * 10);

  for (let i = 0; i < 20 + randomNumber; i++) {
    randomString += String.fromCharCode(33 + Math.floor(Math.random() * 94));
  }

  return randomString;
}

function getCookie(name) {
  let dc = document.cookie;
  let prefix = name + "=";
  let begin = dc.indexOf("; " + prefix);
  if (begin === -1) {
    begin = dc.indexOf(prefix);
    if (begin !== 0) return null;
  }
  else
  {
    begin += 2;
    var end = document.cookie.indexOf(";", begin);
    if (end === -1) {
      end = dc.length;
    }
  }
  // because unescape has been deprecated, replaced with decodeURI
  //return unescape(dc.substring(begin + prefix.length, end));
  return decodeURI(dc.substring(begin + prefix.length, end));
}

export default {
  name: "MainPage",
  components: {LogoImport},
  setup() {
    return {
      DiscordLogo,
      HellasAIOLogo
    };
  },
  created() {
    if(getCookie("accessBool") !== null) {
      router.push('authenticated')
    }
  },
  methods: {
    redirectToDiscord() {
      const randomString = generateRandomString();
      localStorage.setItem('oauth-state', randomString);
      this.$refs.discordButton.setAttribute('disabled', '');
      window.location.href = 'https://discord.com/api/oauth2/authorize?client_id=992931841722032168&redirect_uri=https%3A%2F%2Fquicktask.hellasaio.com%2Fkey&response_type=code&scope=identify' + `&state=${btoa(randomString)}`;
    }
  }
}
</script>