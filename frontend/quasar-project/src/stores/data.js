import { defineStore } from 'pinia'
import { LocalStorage } from 'quasar'
export const dataStore = defineStore('data', {
  state: () => ({
    folderNames: '',
    folderIds: '',
    edit: {
      id: '',
      folderId: '',
      url: '',
      title: ''
    },
    isShowTrashedBookmark: false
  }),
  getters: {
    getFolders: () => {
      var folderNames = LocalStorage.getItem('folderNames').toString()
      return folderNames.split('|').filter(Boolean)
    },
    getFolderId: () => {
      return (folderName) => {
        if (folderName === '') { return '' }
        return LocalStorage.getItem(folderName)
      }
    },
    getFolderName: () => {
      return (folderId) => {
        return LocalStorage.getKey(folderId).toString()
      }
    },
    getEditTargetData (state) {
      return state.edit
    },
    isShowTrashed (state) { return state.isShowTrashedBookmark }
  },
  actions: {
    setFolders (foldersList) {
      console.log(foldersList)
      foldersList.forEach(folder => {
        this.folderNames += folder.FolderName + '|'
        this.folderIds += folder.FolderID + '|'
        LocalStorage.set(folder.FolderName, folder.FolderID)
      })
      LocalStorage.set('folderNames', this.folderNames)
    },
    setEditTargetData (id, title, url, folderid) {
      this.edit.id = id
      this.edit.title = title
      this.edit.url = url
      this.edit.folderId = folderid
    },
    setIsShowTrashed (flagValue) {
      this.isShowTrashedBookmark = flagValue
    }
  }
})
