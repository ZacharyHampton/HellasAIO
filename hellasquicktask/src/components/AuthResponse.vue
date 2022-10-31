<template>
  <br><br><br>
  <div v-if="success">
    <p class="font-bold text-green-600">{{ message }}</p>
  </div>

  <div v-if="success === false">
    <p class="font-bold text-red-600">{{ message }}</p>
  </div>
</template>

<script>
export default {
  name: "AuthResponse",
  props: ['success', 'message']
}
</script>