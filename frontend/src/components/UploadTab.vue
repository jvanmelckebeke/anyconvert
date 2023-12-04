<template>
  <div class="flex">
    <h2>To JPG</h2>
    <div class="grid-container ">
      <div class="grid-item">
        <FileUpload mode="basic" name="file" :url="uploadUrl"
                    auto
                    :accept="accept"
                    @upload="onUpload"></FileUpload>
      </div>

      <br>

      <div class="grid-item">
        <slot name="result" :resultUrl="resultUrl">
          <p v-if="resultUrl">
            {{ resultUrl }} this is default
          </p>
        </slot>
      </div>
    </div>
  </div>
</template>

<script>

import FileUpload from "primevue/fileupload";
import Divider from "primevue/divider";

export default {
  components: {
    FileUpload,
    Divider
  },
  props: {
    accept: String,
    endpoint: String,
    uploadPath: String

  },
  data() {
    return {
      resultUrl: null,
      uploadUrl: `${this.endpoint}${this.uploadPath}`
    };
  },
  methods: {
    onUpload(event) {
      const xhr = event.xhr;
      console.log("debug", xhr)

      // extract response from xhr
      const response = JSON.parse(xhr.response);

      console.log(response)

      this.resultUrl = this.endpoint + response.file;

    },
  },
};
</script>

<style>
.grid-container {
  //display: grid;
  //grid-template-columns: 45% 10% 45%;
  //grid-template-rows: 1fr;
  //margin: 1rem;
}
</style>