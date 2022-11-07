<template>
  <q-page class="bg-dark window-height window-width row justify-center items-center">
    <div class="column">
      <div class="row">
        <q-card square bordered class="q-pa-lg shadow-2">
          <h5 class="text-h5 text-blue text-center q-my-lg">Bookmark Management Tool</h5>
          <q-card-section>
            <q-form class="q-gutter-md">
              <q-input square clearable v-model="userID" type="text" label="userID" />
              <q-input square clearable v-model="password" type="password" label="password" />
            </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
            <q-btn unelevated color="light-blue-8" size="lg" class="full-width" label="Login" @click="login()" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none" @click="gotoRegistUserPage()">
            <p class="text-grey-6 gotoRegistPage">Not reigistered? Created an Account</p>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script>
import { defineComponent } from 'vue'
import { api } from 'boot/axios'
import { useAuthStore } from 'stores/auth'
import { useQuasar } from 'quasar'
export default defineComponent({
  name: 'LoginCard',
  data () {
    return {
      userID: '',
      password: ''
    }
  },
  methods: {
    login () {
      const store = useAuthStore()
      api.post('/login', {
        userid: this.userID,
        password: this.password
      }).then(response => {
        if (response.data.code === 200) {
          store.setToken(response.data.token, response.data.expire)
          window.setTimeout(() => {
            this.$router.push('/app')
          }, 3000)
        }
      }).catch(() => {
        const $q = useQuasar()
        $q.notify({
          message: 'failed login user',
          color: 'negative'
        })
      })
    },
    gotoRegistUserPage () {
      this.$router.push('/regist')
    }
  }
})
</script>

<style scoped>
.q-card {
  width: 360px;
}
.q-input {
    font-size: 18px;
}
.gotoRegistPage:hover{
  font-weight: bold;
  cursor: pointer;
}
</style>
