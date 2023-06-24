<script lang="ts">
  import Dropzone from "svelte-file-dropzone/Dropzone.svelte";
  import { writable } from "svelte/store";
  import { Canvas, Layer } from "svelte-canvas";
  import type { RenderFunc } from "$lib/types";

  const templateFile = writable<HTMLImageElement>();
  let hasFile = false;

  function handleFilesSelect(e: any) {
    if ('detail' in e && 'acceptedFiles' in e.detail) {
      hasFile = true;
      const f: File = e.detail.acceptedFiles[0];
      f.arrayBuffer().then((t) => {
        const img = new Image();
        img.src = "data:image/webp;base64," + btoa(String.fromCharCode(...new Uint8Array(t)));
        img.onload = () => templateFile.set(img);
      })
    }
  }

  $: render = (({ context }) => {
    if ($templateFile) {
      context.drawImage($templateFile, 0, 0)
    }
  }) as RenderFunc;
</script>

<div class="max-w-screen-2xl">
  {#if !hasFile}
  <Dropzone on:drop={handleFilesSelect} multiple={false} accept="image/*" required={true} disableDefaultStyles={true} containerClasses="dropper">
    Drag 'n' drop a meme template, or click to select
  </Dropzone>
  {:else}
    <Canvas width={640} height={640}>
      <Layer {render} />
    </Canvas>
  {/if}
</div>