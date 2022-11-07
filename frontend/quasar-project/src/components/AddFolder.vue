<template>
  <div>
    <q-btn outline rounded class="addFolderBtn"
     data-form="tes" @click="addFolderForm = true"
     >
      <i class='fa fa-plus' data-form="tes"></i>
    </q-btn>
    <q-dialog
      v-model="addFolderForm"
    >
      <folder-name-input :isRename="false"></folder-name-input>
    </q-dialog>
  </div>
</template>
<script>
import { ref } from 'vue'
import FolderNameInput from './FolderNameInput.vue'
export default {
  components: {
    FolderNameInput
  },
  setup () {
    return {
      addFolderForm: ref(false)
    }
  }
}
</script>
<style scoped>

@import url(//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css);

@import url(https://fonts.googleapis.com/css?family=Titillium+Web:300);
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
