<template>
  <div class="flex items-center justify-center h-screen font-sans flex-col">
    <LogoImport title="You are already authenticated. Have fun!"/>
  </div>
</template>

<script>
import LogoImport from "@/components/LogoImport";
export default {
  name: "AlreadyAuthenticated",
  components: {LogoImport}
}
</script>