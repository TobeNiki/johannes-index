<template>
<div class="q-pa-xs" style="max-width: 500px">
  <q-list bordered>
    <q-virtual-scroll class="folderScroll"
        style="max-height: 50vh;"
        :items="folderList"
        separator
        v-slot="{ item, index }"
      >
    <q-item clickable v-ripple class="folder"  :key="index" @click="$emit('eventFolder', item.FolderID)">
      <q-item-section avatar>
        <q-avatar rounded>
          <i class="fa fa-folder "/>
        </q-avatar>
      </q-item-section>
      <q-item-section>{{item.FolderName}}</q-item-section>
      <context-menu :isDelete="true"  :isFolder="true" :folderId="item.FolderID" :targetName="`${item.FolderID}`"/>
    </q-item>
    </q-virtual-scroll>
  </q-list>
</div>
</template>
<script>
import { ref } from 'vue'
import { api } from 'boot/axios'
import { useAuthStore } from 'stores/auth'
import { dataStore } from 'stores/data'
import { useQuasar } from 'quasar'
import ContextMenu from '../components/ContextMenu.vue'
export default {
  name: 'FolderList',
  components: {
    ContextMenu
  },
  setup () {
    const store = useAuthStore()
    const storeData = dataStore()
    const folderList = ref(null)
    const $q = useQuasar()
    api.get('/auth/folder', {
      headers: {
        Authorization: store.getAuthrizationHeader
      }
    })
      .then(response => {
        if (response.data.result.length > 0) {
          folderList.value = response.data.result
          storeData.setFolders(folderList.value)
        }
      })
      .catch(() => {
        $q.notify({
          message: 'failed load folder',
          color: 'negative'
        })
      })
    return {
      folderList
    }
  }
}
</script>
<style scoped>
.folder {
    margin: 10px;
    border-radius: 5px;
    border: 2px solid #313131;
}
.folderScroll::-webkit-scrollbar{
  display: none;
}
</style>
