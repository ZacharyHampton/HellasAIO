<template>
  <div class="flex items-center justify-center h-screen font-sans flex-col">
    <LogoImport class="mb-48"/>
    <AuthResponse :success=success :message=message />
  </div>
</template>

<script>
import AuthResponse from "@/components/AuthResponse";
import LogoImport from "@/components/LogoImport";
import axios from "axios";
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
  name: "QuicktaskPage",
  components: {AuthResponse, LogoImport},
  data() {
    return {
      success: null,
      message: ""
    }
  },
  created() {
    if(getCookie("accessBool") === null) {
      router.push('/')
    }

    const fragment = new URLSearchParams(window.location.search);
    const [siteId, MSKU, size] = [fragment.get('siteId'), fragment.get('product_id'), fragment.get('size')]

    const data = {
      siteId: siteId,
      product_id: MSKU,
      size: size
    }

    const params = new URLSearchParams(data);

    axios.defaults.withCredentials = true
    axios.get('https://api.hellasaio.com/api/quicktask/start', { params })
        .then((response) => {
          this.success = response.status === 200;
          try {
            this.message = response.data.message
          } catch (e) {
            this.message = response.data
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
  }
}
</script>