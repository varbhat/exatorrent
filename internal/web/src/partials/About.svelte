<script lang="ts">
  import { onMount } from 'svelte';
  import { versionstr, versionchecked, Send } from './core';

  onMount(() => {
    if ($versionchecked === false) {
      Send({
        command: 'version'
      });
    }
  });
</script>

<div class="mx-auto max-w-3xl text-gray-200 mt-3">
  <p class="text-center text-lg mb-1 text-blue-300">exatorrent is torrent client</p>
  <p class="text-center mb-2 font-mono">License: GPLv3</p>
  <p class="text-center mb-5 font-mono">Version: {$versionstr}</p>
  <div class="flex justify-around">
    <a href="https://github.com/varbhat/exatorrent" target="_blank" rel="noopener noreferrer">
      <button type="button" class="my-2 mx-2 flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-secGray-700 no-HL font-mono focus:outline-none">Source</button>
    </a>
  </div>
</div>
