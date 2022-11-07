<template>
    <q-page class="bg-dark window-height window-width row justify-center items-center">
      <q-card dark flat bordered class="editCard" >
        <q-card-section>
          <div class="text-h5" style="text-align: center;">Edit : {{ title }}</div>
        </q-card-section>
        <q-card-section>
          {{ text }}
        </q-card-section>
        <q-separator dark inset />
        <q-card-section >
          <div class="text-h6 q-ml-md"> Folder </div>
          <q-select v-model="folder" class="q-ma-md" dark filled label="Folder" :options="options" />
        </q-card-section>
        <q-card-section>
          <div class="text-h6 q-ml-md"> Url </div>
          <q-input v-model="url" class="q-ma-md" dark dense label="URL" />
        </q-card-section>
        <div class="btnGroup" >
          <q-btn flat class="q-ma-lg updateBtn" color="primary" label="Update" @click="updateBookmark()"/>
          <q-btn flat class="q-ma-lg deleteBtn" label="Delete" @click="deleteBookmark()" />
        </div>
      </q-card>
    </q-page>
</template>
<script>
import { useAuthStore } from 'stores/auth'
import { api } from 'boot/axios'
import { defineComponent, ref } from 'vue'
import { dataStore } from 'stores/data'
import { useQuasar } from 'quasar'
export default defineComponent({
  name: 'EditBookmarkPage',
  setup () {
    const storeData = dataStore()
    const editTargetData = storeData.getEditTargetData
    const title = ref(editTargetData.title)
    const text = ref(editTargetData.text)
    const folder = ref('')
    const url = ref(editTargetData.url)
    const id = editTargetData.id
    const $q = useQuasar()
    const store = useAuthStore()
    // load bookmark from id
    return {
      id,
      title,
      text,
      folder,
      url,
      options: storeData.getFolders,
      updateBookmark () {
        api.put('/auth/bookmark', {
          id: this.id,
          url: this.url,
          isUseChromeDriver: false,
          folderId: storeData.getFolderId(this.folder)
        }, {
          headers: {
            Authorization: store.getAuthrizationHeader
          }
        })
          .then(response => {
            $q.notify({
              message: response.data.message,
              color: 'positive'
            })
            this.$router.push('/app')
          })
          .catch(() => {
            $q.notify({
              message: 'failed update bookmark',
              color: 'negative'
            })
          })
      },
      deleteBookmark () {
        api.delete('/auth/bookmark', {
          params: {
            bookmarkid: this.id
          },
          headers: {
            Authorization: store.getAuthrizationHeader
          }
        }).then((response) => {
          $q.notify({
            message: response.data.message,
            color: 'positive'
          })
          this.$router.push('/app')
        })
          .catch(() => {
            $q.notify({
              message: 'failed delete bookmark',
              color: 'negative'
            })
          })
      }
    }
  }
})
</script>
<style scoped>
.editCard {
    width: 500px;
}
.btnGroup {
    float: right;
}
.updateBtn {
    border: solid 1px;
}
.deleteBtn {
    border: solid 1px #ff0070;
    color: #ff0070
}
</style>
