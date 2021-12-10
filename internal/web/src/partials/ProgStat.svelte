<script lang="ts">
  import { fileSize } from './core';

  export let bytescompleted: number;
  export let length: number;
  export let offset: number;

  let progress: number;

  $: progress = (bytescompleted / length) * 100;
</script>

<div class="space-y-2 mt-1">
  <div class="bg-black rounded-full overflow-hidden">
    <div class="bg-blue-700  h-1.5" style="width:{progress ? progress : 0}%" />
  </div>
  <div class="text-neutral-400 flex justify-between text-sm font-medium">
    <div class="break-all">
      {fileSize(bytescompleted)} / {fileSize(length)} (Off. {offset})
    </div>
    <div>
      {progress?.toLocaleString('en-US', {
        maximumFractionDigits: 2,
        minimumFractionDigits: 2
      })} %
    </div>
  </div>
</div>
