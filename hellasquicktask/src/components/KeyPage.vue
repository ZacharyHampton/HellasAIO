<template>
<div class="flex items-center justify-center h-screen font-sans flex-col">
  <LogoImport title="" class="mb-48"></LogoImport>

  <p class="font-semibold">Bot Key:</p>
  <input type="text" placeholder="Bot Key" class="border-2 border-gray-300 hover:border-black rounded-lg text-sm mb-2" ref="key">
  <button @click="login()" ref="submitButton" class="flex box-border items-center text-white bg-blue-600 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 disabled:bg-slate-50 disabled:text-slate-500 disabled:border-slate-200 font-medium rounded-lg text-sm px-5 py-2.5 mr-2 mb-2">Submit</button>

  <AuthResponse :success=success :message=message ref="aResponse"/>
</div>
</template>

<script>
import HellasAIOLogo from "@/assets/logo.png";
import LogoImport from "@/components/LogoImport";
import axios from 'axios';
import AuthResponse from "@/components/AuthResponse";
import router from "@/router";

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
  name: "KeyPage",
  components: {AuthResponse, LogoImport},
  setup() {
    return {
      HellasAIOLogo
    };
  },
  data() {
    return {
      code: null,
      message: "",
      success: null,
    }
  },
  created() {
    if(getCookie("accessBool") !== null) {
      router.push('authenticated')
    }

    const fragment = new URLSearchParams(window.location.search);
    const state = fragment.get('state');
    this.code = fragment.get('code');

    if (localStorage.getItem('oauth-state') !== atob(decodeURIComponent(state))) {
      return console.log('You may have been clickjacked!'); // make this proper
    }
  },
  methods: {
    login() {
      this.$refs.submitButton.setAttribute('disabled', '');
      axios.defaults.withCredentials = true
      axios.post('https://api.hellasaio.com/api/quicktask/auth', {'key': this.$refs.key.value, 'code': this.code})
          .then((response) => {
            if(response.status !== 200) {
              this.success = false
              try {
                this.message = response.data.message
              } catch (e) {
                this.message = response.data
              }
            } else {
              this.success = true
              this.message = "Successfully authenticated."
            }
          })
          .catch((error) => {
            console.log(error);
            this.success = false
            try {
              this.message = error.response.data.message
            } catch (e) {
              this.message = "Request error."
            }
          })

      this.$refs.submitButton.removeAttribute('disabled')
    }
  }
}
</script>