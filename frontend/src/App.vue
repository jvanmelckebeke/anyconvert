<template>
  <div id="app">
    <h1>File Converter</h1>
    <TabView>
      <TabPanel header="toJpg">
        <div>
          <UploadTab
              endpoint="http://localhost:8000"
              accept="image/webp"
              upload-path="/image"
          >
            <template #result="result">
                <Image :src="result.resultUrl" image-style="max-width: 100%">

                </Image>
            </template>
          </UploadTab>
        </div>
      </TabPanel>
      <TabPanel header="toMp4">
        <div>
          <UploadTab
              endpoint="http://localhost:8000"
              accept="image/webp"
              upload-path="/image">
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
