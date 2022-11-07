<template>
     <q-menu
        touch-position
        context-menu
      >
        <q-list dense style="min-width: 100px">
          <q-item clickable v-close-popup @click="ShowCount(targetName)">
            <q-item-section>Count</q-item-section>
          </q-item>
          <q-item clickable @click="renameFolderForm = true" v-if="isFolder">
            <q-item-section>Rename</q-item-section>
          </q-item>
            <q-dialog v-model="renameFolderForm" >
              <folder-name-input :isRename="true" :folderID="folderId"></folder-name-input>
            </q-dialog>
          <q-separator />
          <q-item clickable v-close-popup v-if="isDelete" @click="Delete(isFolder, folderId)">
            <q-item-section>delete</q-item-section>
          </q-item>
        </q-list>

      </q-menu>
</template>
<script>
import { ref } from 'vue'
import FolderNameInput from './FolderNameInput.vue'
import { useQuasar } from 'quasar'
import { api } from 'boot/axios'
import { useAuthStore } from 'stores/auth'
export default {
  components: {
    FolderNameInput
  },
  props: {
    targetName: {
      type: String,
      required: true
    },
    isFolder: {
      type: Boolean,
      required: true
    },
    folderId: {
      type: String
    },
    isDelete: {
      type: Boolean,
      required: true
    }
  },
  setup () {
    const store = useAuthStore()
    const $q = useQuasar()
    const renameFolderForm = ref(false)
    return {
      renameFolderForm,
      ShowCount (targetName) {
        api.get('/auth/bookmark/count', {
          params: {
            target: targetName
          },
          headers: {
            Authorization: store.getAuthrizationHeader
          }
        })
          .then(response => {
            $q.notify({
              position: 'top',
              message: response.data.result + ' Bookmark',
              multiLine: true,
              type: 'info'
            })
          })
          .catch(err => {
            console.log(err)
            $q.notify({
              message: 'failed load propaty',
              color: 'negative'
            })
          })
      },
      Delete (isFolder, folderId) {
        if (isFolder) {
          api.delete('/auth/folder', {
            params: {
              folderid: folderId
            },
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
                message: 'failed delete folder',
                color: 'negative'
              })
            })
        }
      }
    }
  }
}
</script>
