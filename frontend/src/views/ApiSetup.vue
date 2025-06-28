<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>Typhoon API Configuration</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <p class="mb-4">
              Please enter the API base URL and Token for your Typhoon backend.
            </p>
            <v-form @submit.prevent="submitConfig">
              <v-text-field
                label="API Base URL"
                v-model="apiUrl"
                prepend-icon="mdi-web"
                required
                :rules="[v => !!v || 'API URL is required']"
                placeholder="e.g., http://localhost:8080/api/v1"
              ></v-text-field>
              <v-text-field
                label="API Token (Optional)"
                v-model="apiToken"
                prepend-icon="mdi-key-variant"
                type="password"
                hint="Leave blank if your backend does not require a token."
              ></v-text-field>
              <v-alert v-if="errorMessage" type="error" dense class="mt-3 mb-3">
                {{ errorMessage }}
              </v-alert>
              <v-btn type="submit" color="primary" block large class="mt-4" :loading="loading">
                Save and Connect
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import apiConfig, { saveApiConfig } from '@/services/apiConfig';
import { checkBackendConnection } from '@/services/api'; // Assuming api.js will have this

const router = useRouter();
const apiUrl = ref(apiConfig.apiUrl || '');
const apiToken = ref(apiConfig.apiToken || '');
const errorMessage = ref('');
const loading = ref(false);

async function submitConfig() {
  if (!apiUrl.value) {
    errorMessage.value = 'API URL cannot be empty.';
    return;
  }
  loading.value = true;
  errorMessage.value = '';

  // Temporarily save config to be used by checkBackendConnection
  saveApiConfig(apiUrl.value, apiToken.value);

  try {
    // Attempt to connect to a basic backend endpoint to verify
    // For example, fetching version or a dedicated health check endpoint
    // This function needs to be implemented in api.js
    await checkBackendConnection();

    // If successful, the config is already saved by the watcher in apiConfig.js
    // Navigate to the dashboard or home page
    router.push('/');
  } catch (error) {
    errorMessage.value = `Failed to connect to backend: ${error.message}. Please check the URL and token.`;
    // Optionally clear the temporarily saved config if connection fails
    // clearApiConfig(); // Or let user retry with current values
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.fill-height {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
