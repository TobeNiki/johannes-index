<template>
  <q-page class="flex flex-center">
    <div class="folderTree" >
  <div class="q-pa-md bg-grey-10 text-white">
    <q-list dark  separator >
      <q-item clickable v-ripple class="folderListItem" @click.left="showBookmark('/auth/bookmark')">
        <i class="fa fa-2x fa-book"></i>
        <q-item-section>
          <span>All Bookmark</span>
        </q-item-section>
        <context-menu :isDelete="false"  :isFolder="false" :targetName="'all'"/>
      </q-item>

      <q-item clickable v-ripple class="folderListItem" @click.left="showBookmark('/auth/bookmark/unorganized')">
        <i class="fa fa-box fa-2x"></i>
        <q-item-section>
          <span>Unorganized</span>
        </q-item-section>
        <context-menu :isDelete="false" :isFolder="false" :targetName="'unorganized'"/>
      </q-item>

      <q-item clickable v-ripple class="folderListItem" @click.left="showBookmark('auth/bookmark/trash')">
          <i class="fa fa-trash-o fa-2x"></i>
          <q-item-section>
            <span>Trash</span>
          </q-item-section>
          <context-menu :isDelete="true" :isFolder="false" :targetName="'trashed'"/>
      </q-item>
    </q-list>
  </div>
  <div class="hr-line"></div>
  <div class="folderTreeLabel" >
    Folder
    <add-folder></add-folder>
  </div>
  <div class="q-pa-none bg-grey-10 text-white">
    <folder-list @eventFolder="getBookmarkFromFolder"/>
  </div>
  <div class="hr-line"></div>
  <div class="accountSetting">
    <div class="q-pa-md bg-grey-10 text-white">
      <q-list dark padding bordered class="rounded-borders" style="max-width: 328px">
      <q-expansion-item
        icon="perm_identity"
        label="Account settings"
      >
      </q-expansion-item>
      </q-list>
    </div>
  </div>
  </div>
    <search-bar @eventSearch="searchRun"/>
    <add-button />
    <bookmark-table :bookmarks="bookmarks"/>
  </q-page>
</template>

<script>

import { ref } from 'vue'
import FolderList from '../components/FolderList.vue'
import SearchBar from '../components/SeachBar.vue'
import AddButton from '../components/AddButton.vue'
import BookmarkTable from '../components/BookmarkTable.vue'
import AddFolder from '../components/AddFolder.vue'
import ContextMenu from '../components/ContextMenu.vue'
import { api } from 'boot/axios'
import { bookmarkStore } from 'stores/bookmark'
import { useAuthStore } from 'stores/auth'
import { dataStore } from 'stores/data'
import { setBookmarkData } from '../module/api.js'
import { useQuasar } from 'quasar'
export default {
  name: 'IndexPage',
  components: {
    FolderList,
    BookmarkTable,
    SearchBar,
    AddButton,
    AddFolder,
    ContextMenu
  },
  setup () {
    const $q = useQuasar()
    const bookmarks = ref(null)
    const storeData = dataStore()
    const store = useAuthStore()
    const bmStore = bookmarkStore()
    const isActive = ref(false)
    const showBookmark = (endpoint) => {
      if (endpoint === 'auth/bookmark/trash') {
        storeData.setIsShowTrashed(true)
      }
      api.get(endpoint, {
        params: {
          top: 10,
          sort: bmStore.getSortType
        },
        headers: {
          Authorization: store.getAuthrizationHeader
        }
      })
        .then(response => setBookmarkData(response, bookmarks))
        .catch(() => {
          $q.notify({
            message: 'failed load bookmark',
            color: 'negative'
          })
        })
    }
    const searchRun = (searchWord) => {
      if (searchWord === '') { return }
      storeData.setIsShowTrashed(false)
      api.post('/auth/bookmark/search', {
        word: searchWord,
        search_word_type_is_and: false,
        search_target: 'text',
        top: 10,
        folderid: '',
        sort_type: bmStore.getSortType
      }, {
        headers: {
          Authorization: store.getAuthrizationHeader
        }
      })
        .then(response => setBookmarkData(response, bookmarks))
        .catch(() => {
          $q.notify({
            message: 'failed search bookmark',
            color: 'negative'
          })
        })
    }
    const getBookmarkFromFolder = (folderId) => {
      storeData.setIsShowTrashed(false)
      api.get('/auth/bookmark/folder', {
        params: {
          top: 10,
          folderid: folderId,
          sort: bmStore.getSortType
        },
        headers: {
          Authorization: store.getAuthrizationHeader
        }
      })
        .then(response => setBookmarkData(response, bookmarks))
        .catch(() => {
          $q.notify({
            message: 'failed get bookmark from folder',
            color: 'negative'
          })
        })
    }
    showBookmark('/auth/bookmark')
    return {
      isActive,
      bookmarks,
      showBookmark,
      searchRun,
      getBookmarkFromFolder
    }
  }
}

</script>
<style scoped>

@import url(//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css);

@import url(https://fonts.googleapis.com/css?family=Titillium+Web:300);
.hr-line {
    border-top: 3px solid #313131;
    margin: 0;
}
.folderListItem{
    text-align: left;
    padding: 15px 20px;
}
.folderListItem i {
    width: 23px;
    margin-right: 15px;
}
.folderListItem span {
    font-size: 16px;
    font-family: "Roboto", "-apple-system", "Helvetica Neue", Helvetica, Arial, sans-serif;
}
.search{
    margin: 10px;
}
.folderTree {
    width: 30vh;
    position: absolute;
    left:0;
    top: 0;
    background-color: #212121;
    min-height: calc(100vh - 50px);
}
.folderTreeLabel {
    margin: 20px 5px 0px 5px;
    padding-left: 70px;
    font-size: 18px;
    color: #fff;
    font-family: "Roboto", "-apple-system", "Helvetica Neue", Helvetica, Arial, sans-serif;
}
.folderTreeLabel  {
    margin-left: 10px;
    min-width: 7vh;
}
.newFolderForm {
    background-color: #313131;
    border-radius: 10px;
    font-size: 18px;
    padding-left: 20px;
    font-family: "Roboto", "-apple-system", "Helvetica Neue", Helvetica, Arial, sans-serif;
}
.flex-center {
  background-color: #353535;
}
.addFolderBtn {
    border: solid 2px #313131;
    min-width: 7vh;
    margin-left: 80px;
    margin-top: -40px;
}

.addFolderBtn:hover {
  transition: 1200ms;
  border: solid 2px #1976d2;
}
</style>
