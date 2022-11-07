import { defineStore } from 'pinia'
export const bookmarkStore = defineStore('bookmark', {
  state: () => ({
    sortType: 'desc'
  }),
  getters: {
    getSortType: (state) => state.sortType
  },
  actions: {
    setSortType (sort) {
      this.sortType = sort
    }
  },
  persist: {
    enabled: true,
    strategies: [
      { storage: localStorage }
    ]
  }
})
