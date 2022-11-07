<template>
  <div class="bookMarkTable">
    <div >
      <q-virtual-scroll class="bookmarkscroll"
        style="max-height: 85vh;"
        :items="bookmarks"
        separator
        v-slot="{ item, index }"
      >
        <div class="my-card" :key="index">
          <q-card flat dark bordered >
            <q-item>
              <q-item-section avatar >
                <q-avatar size="16px">
                  <img :src="`data:image/png;base64,${item.favicon}`" alt="icon">
                </q-avatar>
              </q-item-section>
              <q-item-section>
                <q-item-label>
                  <div class="text-h6 q-mt-sm q-mb-xs">{{ item.title }}</div>
                </q-item-label>
                <div class="caption">
                  <q-item-label caption  class="label">
                    {{ item.date }}
                  </q-item-label>
                  <q-item-label class="label">
                    {{ item.url }}
                    <q-popup-edit  v-model="item.url" v-slot="scope">
                      <q-input v-model="scope.value" dense autofocus @keyup.enter="urlEdit(scope.value)" />
                    </q-popup-edit>
                  </q-item-label>
                </div>
                <div class="btnGroup">
                <q-btn @click="tolinkGo(item.url)" class="toBtn">
                  <i class='fa fa-external-link'></i>
                </q-btn>
                <q-btn class="toBtn" @click="editBookmark(item.id, item.title, item.url, item.folderId)">
                  <i class="fa fa-pencil-square-o"></i>
                </q-btn>
                <q-btn @click="toTrash(item.id)" class="toBtn">
                  <i class="fa fa-trash-o"></i>
                </q-btn>
                </div>
              </q-item-section>
            </q-item>
            <q-separator class="separator" />
            <q-card-section horizontal >
              <q-card-section>
                {{ item.text }}
              </q-card-section>
            </q-card-section>
          </q-card>
        </div>
      </q-virtual-scroll>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { api } from 'boot/axios'
import { dataStore } from 'stores/data'
import { useAuthStore } from 'stores/auth'
import { useQuasar } from 'quasar'
import { sanitizeUrl } from '@braintree/sanitize-url'
export default {
  props: {
    bookmarks: Object
  },
  setup () {
    const storeData = dataStore()
    const store = useAuthStore()
    const $q = useQuasar()
    return {
      isEdit: ref(false),
      options: storeData.getFolders,
      newFolder: ref(''),
      toTrash (id) {
        if (id === '') { return }
        if (storeData.isShowTrashed) {
          api.delete('/auth/bookmark', {
            params: {
              bookmarkid: id
            },
            headers: {
              Authorization: store.getAuthrizationHeader
            }
          }).then((response) => {
            $q.notify({
              message: response.data.message,
              color: 'positive'
            })
          })
            .catch(() => {
              $q.notify({
                message: 'failed delete bookmark',
                color: 'negative'
              })
            })
          return
        }
        api.put('/auth/bookmark/trash', {}, {
          params: { bookmarkid: id },
          headers: {
            Authorization: store.getAuthrizationHeader
          }
        }).then(response => {
          $q.notify({
            message: response.data.message,
            color: 'positive'
          })
        }).catch(() => {
          $q.notify({
            message: 'failed trash bookmark',
            color: 'negative'
          })
        })
      }
    }
  },
  methods: {
    editBookmark (id, title, url, folderId) {
      console.log(id)
      const storeData = dataStore()
      storeData.setEditTargetData(id, title, url, folderId)
      this.$router.push('/edit')
    },
    tolinkGo (url) {
      window.open(sanitizeUrl(url), '_blank', 'noreferrer')
    },
    urlEdit (newUrl) {
      console.log(newUrl)
    }
  },
  watch: {
    bookmarks: function (newBookmark, oldBookmark) {
      console.log(newBookmark, oldBookmark)
    }
  }
}
</script>
<style scoped>
.bookMarkTable {
    min-width: calc(100% - 30vh);
    position: absolute;
    left:30vh;
    top: 6.2vh;
    margin: 0px 0px;
    padding: 1px 0px;
}
.my-card {
  margin: 0px 0px;
}
.my-card :hover{
  background-color: #313131;
}
.separator{
    background-color: #515151;
}
.caption {
  font-size: 14px;
}
.label {
  color: whitesmoke;
  display: inline-block;
  margin: 5px;
}
.btnGroup {
  position: absolute;
  right: 20px;
  top: 15px;
  display: flex;
}
.toBtn {
  border: solid 1.5px #515151;
  border-radius: 5px;
  padding: 0;
  position: relative;
  width: 60px;
  margin: 0 3px;
  background-color:#313131;
}
.toBtn :hover {
  background-color:#515151;
  border: none;
  border-radius: 5px;
}
.bookmarkscroll::-webkit-scrollbar {
  display: none;
}
</style>
