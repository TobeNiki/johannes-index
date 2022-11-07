<template>
      <q-card dark style="width: 400px">
        <q-card-section>
          <div class="text-h6">FolderName</div>
        </q-card-section>

        <q-input class="folderNameInput" dark filled type="text" v-model="newFolderName" :disable="disable">
        </q-input>
        <q-card-actions align="right" class="bg-dark text-teal">
          <q-btn flat label="OK" v-close-popup  @click="NewFolder(isRename, folderID)"/>
        </q-card-actions>
      </q-card>
</template>

<script>
import { ref, defineComponent } from 'vue'
import { api } from 'boot/axios'
import { useAuthStore } from 'stores/auth'
export default defineComponent({
  props: {
    isRename: {
      type: Boolean,
      required: true
    },
    folderID: {
      type: String
    }
  },
  setup () {
    const store = useAuthStore()
    return {
      newFolderName: ref(''),
      NewFolder (isRename, folderID) {
        if (isRename) {
          api.put('/auth/folder/rename', {
            folderid: folderID,
            foldername: this.newFolderName
          }, {
            headers: {
              Authorization: store.getAuthrizationHeader
            }
          })
            .then(() => {
              this.newFolderName = ''
              this.$router.go(0)
            })
            .catch(err => {
              alert(err)
            })
        } else {
          api.post('/auth/folder', {
            foldername: this.newFolderName
          }, {
            headers: {
              Authorization: store.getAuthrizationHeader
            }
          })
            .then(() => {
              this.newFolderName = ''
              this.$router.go(0)
            })
            .catch(err => {
              alert(err)
            })
        }
      }
    }
  }
})
</script>
<style scoped>
.folderNameInput {
    margin: 0px 25px;
}
</style>
