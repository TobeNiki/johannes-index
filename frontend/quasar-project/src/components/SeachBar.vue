<template>
  <div class="searchBar">
    <q-input dark dense filled v-model="search" class="searchBarInput" :disable="disable" v-on:keydown.enter="$emit(this.search)">
        <template v-slot:append>
          <q-icon name="search" @click="$emit('eventSearch',this.search)"/>
        </template>
    </q-input>
    <sort-dropdown />
  </div>
</template>

<script>
import { ref } from 'vue'
import SortDropdown from './SortDropdown.vue'
export default {
  components: {
    SortDropdown
  },
  setup () {
    return {
      search: ref('')
    }
  }
}
</script>
<style scoped>
.searchBar {
    min-width: calc(100% - 40vh);
    position: absolute;
    left:30vh;
    top: 0vh;
    margin: 0;
    background-color: #171717;
    text-align: left;
    height: 57px;
}
.searchBarInput {
  width: 350px;
  margin: 8px 50px;
}
</style>
