<template>
  <div id="app">
    <h1>File Converter</h1>
    <TabView>
      <TabPanel header="Image">
        <div>
          <UploadTab
              header="To jpg"
              endpoint="http://localhost:8000"
              accept="image/webp"
              upload-path="/image"
          >
            <template #result="result">
              <div class="output">
                <Image :src="result.resultUrl" image-style="max-width: 100%"/>
              </div>
            </template>
          </UploadTab>
        </div>
      </TabPanel>
      <TabPanel header="Video">
        <div>
          <UploadTab
              header="To mp4"
              endpoint="http://localhost:8000"
              accept="image/webp, video/webm, image/gif"
              upload-path="/video">
            <template #result="result">
              <div class="output">

                <video :src="result.resultUrl" controls height="240"/>
              </div>
            </template>
          </UploadTab>
        </div>
      </TabPanel>
    </TabView>
  </div>
</template>

<script>
import UploadTab from "@/components/UploadTab.vue";

import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";

import Image from "primevue/image";

import 'primevue/resources/themes/md-dark-indigo/theme.css'

export default {
  components: {
    UploadTab,
    TabView,
    TabPanel,
    Image
  },
  data() {
    return {
      activeTab: "toJpg",
    };
  },

};
</script>


<style scoped>
.title {

}

.output {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>