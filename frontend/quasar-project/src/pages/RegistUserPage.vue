<template>
  <q-page class="bg-dark window-height window-width row justify-center items-center">
    <div class="column">
      <div class="row">
        <q-card square bordered class="q-pa-lg shadow-2">
          <h5 class="text-h5 text-blue text-center q-my-lg">Create User Account</h5>
          <q-card-section>
            <q-form class="q-gutter-md">
              <q-input square clearable v-model="userID" type="text" label="userID" />
              <q-input square clearable v-model="displayName" type="text" label="displayname" />
              <q-input square clearable v-model="password" type="password" label="password" />
            </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
            <q-btn unelevated color="light-blue-8" size="lg" class="full-width" label="Regist" @click="registUser()" />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script>
import { defineComponent } from 'vue'
import { api } from 'boot/axios'
import { useQuasar } from 'quasar'
export default defineComponent({
  name: 'CreateUserPage',
  data () {
    return {
      userID: '',
      displayname: '',
      password: ''
    }
  },
  methods: {
    registUser () {
      const $q = useQuasar()
      api.post('/regist', {
        displayname: this.displayname,
        logindata: {
          userid: this.userID,
          password: this.password
        }
      }).then(response => {
        if (response.data.code === 201) {
          $q.notify({
            message: 'succes create user',
            color: 'positive'
          })
          setTimeout(() => {
            this.$router.push('/login')
          }, 3000)
        }
      }).catch(() => {
        $q.notify({
          message: 'failed create user',
          color: 'negative'
        })
      })
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
</style>
