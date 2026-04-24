<template>
  <div />
</template>

<script>
export default {
  name: "StripeRedirect",
  mounted() {
    this.$router.replace({ name: "home" })
  },
}
</script>
