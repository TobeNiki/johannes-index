<template>
    <div >
    <q-btn class="addBtn" color="" @click="showAddBookmarkModal()" >
        <i class="fa fa-2x fa-upload" aria-hidden="true"></i>
    </q-btn>
    <q-dialog
      v-model="add"
    >
      <q-card dark style="width: 400px">
        <q-card-section>
          <div class="text-h6">URL</div>
        </q-card-section>
        <q-select v-model="folders" class="urlInput" dark filled label="Folder" :options="options" />
        <q-separator style="margin:5px"/>
        <q-input class="urlInput" dark filled type="url" v-model="url" :disable="disable">
        </q-input>
        <q-card-actions align="right" class="bg-dark text-teal">
          <q-btn flat label="OK" v-close-popup @click="addBookmark()" />
        </q-card-actions>
      </q-card>
    </q-dialog>
    </div>
</template>
<script>
import { ref } from 'vue'
import { api } from 'boot/axios'
import { useAuthStore } from 'stores/auth'
import { dataStore } from 'stores/data'
import { useQuasar } from 'quasar'
export default {
  setup () {
    const storeData = dataStore()
    const $q = useQuasar()
    const store = useAuthStore()
    console.log(storeData.getFolders)
    return {
      add: ref(false),
      url: ref(''),
      folders: ref(''),
      options: storeData.getFolders,
      addBookmark () {
        api.post('/auth/bookmark/add', {
          id: '',
          url: this.url,
          isUseChromeDriver: false,
          folderId: storeData.getFolderId(this.folders)
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
          })
          .catch(() => {
            $q.notify({
              message: 'failed add bookmark',
              color: 'negative'
            })
          })
      },
      showAddBookmarkModal () {
        this.add = true
        this.folders = ''
      }
    }
  }
}
</script>
<style scoped>

@import url(//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css);

@import url(https://fonts.googleapis.com/css?family=Titillium+Web:300);

.addBtn {
    width: 91px;
    height: 56px;
    position: absolute;
    right: 0;
    top: 0vh;
    margin: 0;
    text-align: left;
    background-color: #313131;
    font-size: 16px;
}
.addBtn :hover{
    opacity: 0.5;
}
.urlInput {
    margin: 0px 25px;
}
</style>
