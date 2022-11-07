<template>
    <div class="q-pa-md">
    <q-btn-dropdown  class="sortdd" v-bind:label="sortType">
      <q-list>
        <q-item clickable v-close-popup @click="onItemClick('asc')">
          <q-item-section>
            <q-item-label>Asc</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable v-close-popup @click="onItemClick('desc')">
          <q-item-section>
            <q-item-label>Desc</q-item-label>
          </q-item-section>
        </q-item>

      </q-list>
    </q-btn-dropdown>
  </div>
</template>
<script>
import { bookmarkStore } from 'stores/bookmark'
export default {
  setup () {
    const store = bookmarkStore()
    return {
      sortType: store.getSortType,
      onItemClick (sortType) {
        store.setSortType(sortType)
        this.sortType = sortType
      }
    }
  }
}
</script>
<style scoped>
.sortdd {
    width: 100px;
    position: absolute;
    right: 0;
    top: 0;
    margin: 10px 30px;
    border:solid 2px #313131;
    border-radius: 5px;
    color : gainsboro;
}
.sortdd:hover {
    transition: 1200ms;
    border: solid 2px #1976d2;
}
</style>
