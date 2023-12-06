<template>
  <div class="centered">
    <h2 class="item">{{ header }}</h2>
    <div class="centered-container ">
      <div class="item">
        <FileUpload mode="basic" name="file" :url="uploadUrl"
                    auto
                    :accept="accept"
                    @upload="onUpload"
                    @before-send="onSend"
                    @error="onError"/>
      </div>

      <br>

      <div class="item full-width">
        <slot name="result" :resultUrl="resultUrl" v-if="resultUrl">
          <p v-if="resultUrl">
            {{ resultUrl }} this is default
          </p>
        </slot>

        <slot name="loader" v-if="loading">
          <div class="loader">
            <Skeleton width="100%" height="160px"></Skeleton>
          </div>
        </slot>

        <slot name="error" :errorMessage="errorMessage" v-if="errorMessage">

          <Message severity="error">
            {{ errorMessage }}
          </Message>
        </slot>
      </div>
    </div>
  </div>
</template>

<script>

import FileUpload from "primevue/fileupload";
import Divider from "primevue/divider";
import Message from "primevue/message";
import Skeleton from "primevue/skeleton";


export default {
  components: {
    FileUpload,
    Divider,
    Message,
    Skeleton
  },
  props: {
    header: String,
    accept: String,
    endpoint: String,
    uploadPath: String

  },
  data() {
    return {
      resultUrl: null,
      errorMessage: null,
      loading: false,
      uploadUrl: `${this.endpoint}${this.uploadPath}`
    };
  },
  methods: {
    checkStatus(statusPath) {
      const interval = setInterval(() => {
        fetch(`${this.endpoint}${statusPath}`)
            .then(response => response.json())
            .then(data => {

              if (data.status === "done") {
                this.resultUrl = `${this.endpoint}${data.result_url}`;
                this.loading = false;
                clearInterval(interval);
              } else if (data.status === "error") {
                this.errorMessage = data.error;
                this.loading = false;
                clearInterval(interval);
              }
            })
            .catch(error => {
              console.error(error);
              this.errorMessage = error;
              this.loading = false;
              clearInterval(interval);
            });
      }, 500);
    },
    onUpload(event) {
      const xhr = event.xhr;
      console.log("debug", xhr)

      // extract response from xhr
      const response = JSON.parse(xhr.response);


      // example success response:
      // {
      // "status": "/status/b8724692"
      // }

      this.checkStatus(response.status);

    },
    onSend(event) {
      this.loading = true;
      this.errorMessage = null;
      this.resultUrl = null;
    },
    onError(event) {
      const {xhr, files} = event;

      console.log("error", xhr, files);

      this.errorMessage = xhr.response;
      this.resultUrl = null;
      this.loading = false;

      if (this.errorMessage === "") {
        this.errorMessage = "Something went wrong"
      } else {

        this.errorMessage = xhr.statusText;
      }

    }
  },
};
</script>

<style scoped>
.item {
  text-align: center;
  justify-content: center;
  width: 100%;
  flex-grow:1;
}

.full-width {
  width: 100%;
}

.loader {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>