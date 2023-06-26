<script lang="ts">
  import type { PageData } from "./$types";
  import type { RenderFunc } from "$lib/types";
  import { writable } from "svelte/store";
  import { onMount } from "svelte";
  import { Canvas, Layer } from "svelte-canvas";
  import { HorizontalAlign, Payload, TemplateText, Text, VerticalAlign } from "../../../generated";
  import copy from "copy-to-clipboard";
  import { RadioGroup, RadioItem, SlideToggle, toastStore } from "@skeletonlabs/skeleton";
  import { browser } from "$app/environment";
  import Fa from 'svelte-fa/src/fa.svelte'
  import { faAlignLeft, faAlignCenter, faAlignRight } from '@fortawesome/free-solid-svg-icons';
  import { faCircleDown, faCircleUp, faCircleDot } from '@fortawesome/free-regular-svg-icons';
  import FontFaceObserver from "fontfaceobserver";
  import { PUBLIC_BACKEND_API_BASE } from "$env/static/public";

  const MIN_WIDTH = 50;
  const MIN_HEIGHT = 50;
  const DEFAULT_FONT = 'Impact';
  const DEFAULT_SIZE = 22;
  const DEFAULT_FILL = '#000000';
  const RATIO_MEASURE_FONT_SIZE = 1000;
  const FONTS = [
    'Arial',
    'Courier',
    'Helvetica',
    'Impact',
    'Times New Roman',
  ];
  const URL_REGEX = new RegExp(/[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)/);

  const templateFile = writable<HTMLImageElement>();
  export let data: PageData;

  let loaded = false;
  let hovered = false;
  let cursor = 'initial';

  let width = 0;
  let height = 0;

  let texts: Text[] = [];
  let templateTexts: TemplateText[];

  let fontsReady = writable(false);

  onMount(() => {
    if (!browser) {
      return;
    }

    Promise.any(FONTS.map((f) => new FontFaceObserver(f).load())).then(() => fontsReady.set(true));

    Promise.all([
      fetch(`${PUBLIC_BACKEND_API_BASE}/template/${data.id}.webp`).then(response => response.arrayBuffer()).then((t) => {
        const img = new Image();
        img.src = "data:image/webp;base64," + btoa(String.fromCharCode(...new Uint8Array(t)));
        img.onload = () => {
          width = img.width;
          height = img.height;
          templateFile.set(img);
        };
      }),
      fetch(`${PUBLIC_BACKEND_API_BASE}/api/templates/${data.id}`).then(response => response.json()).then((t) => {
        texts = t.data.map((x: TemplateText, i: number) => {
          return {
            ...TemplateText.fromPartial(x),
            template_text: i
          };
        });
        templateTexts = t.data.map(TemplateText.fromPartial);
      })
    ]).then(() => loaded = true);
  });

  const anchorSize = 5;

  $: render = (({ context }) => {
    if (!$fontsReady) {
      return;
    }

    if ($templateFile) {
      context.drawImage($templateFile, 0, 0);
    }

    [...texts].reverse().forEach(t => {
      if (t.text) {
        const x = t.x || 0;
        const y = t.y || 0;
        const width = t.width || MIN_WIDTH;
        const height = t.height || MIN_HEIGHT;

        let fontSize = (t.size || DEFAULT_SIZE);

        context.fillStyle = t.fill_color || DEFAULT_FILL;
        context.strokeStyle = t.stroke_color || DEFAULT_FILL;
        context.font = fontSize + 'pt ' + (t.font || DEFAULT_FONT);

        let renderX = x;
        switch (t.horizontal_align) {
          default:
          case HorizontalAlign.LEFT:
            context.textAlign = 'left';
            break;
          case HorizontalAlign.CENTER:
            context.textAlign = 'center';
            renderX += width / 2;
            break;
          case HorizontalAlign.RIGHT:
            context.textAlign = 'right';
            renderX += width;
            break;
        }

        let renderY = y;
        switch (t.vertical_align) {
          default:
          case VerticalAlign.TOP:
            context.textBaseline = "top";
            break;
          case VerticalAlign.MIDDLE:
            context.textBaseline = "middle";
            renderY += height / 2;
            break;
          case VerticalAlign.BOTTOM:
            context.textBaseline = "bottom";
            renderY += height;
            break;
        }

        const lines = t.text.trim().split('\n').map(l => l.match(URL_REGEX) ? "Links are not allowed" : l);

        let widest = 0;
        let widestID = 0;
        let lineHeight = 1.35 * fontSize;

        lines.forEach((l, i) => {
          const bb = context.measureText(l);
          if (widest < bb.width) {
            widest = bb.width;
            widestID = i;
          }
        });

        // Check width
        if (widest > width) {
          context.font = RATIO_MEASURE_FONT_SIZE + 'pt ' + (t.font || DEFAULT_FONT);
          const bb = context.measureText(lines[widestID]);
          const ratio = (bb.width - widest) / (RATIO_MEASURE_FONT_SIZE - fontSize);
          fontSize = Math.min(fontSize, Math.floor(width / ratio));
          context.font = fontSize + 'pt ' + (t.font || DEFAULT_FONT);
          lineHeight = 1.35 * fontSize;
        }

        if (lines.length * lineHeight > height) {
          lineHeight = height / lines.length
          fontSize = Math.min(fontSize, lineHeight / 1.35);
          context.font = fontSize + 'pt ' + (t.font || DEFAULT_FONT);
        }

        const totalHeight = (lines.length - 1) * lineHeight;
        switch (t.vertical_align) {
          default:
          case VerticalAlign.TOP:
            // Do nothing
            break;
          case VerticalAlign.MIDDLE:
            renderY -= totalHeight / 2
            break;
          case VerticalAlign.BOTTOM:
            renderY -= totalHeight;
            break;
        }

        lines.forEach((l, i) => {
          if (!t.unfilled) {
            context.fillText(l, renderX, renderY + (i * lineHeight));
          }

          if (t.stroke) {
            context.lineWidth = t.stroke;
            context.strokeText(l, renderX, renderY + (i * lineHeight));
          }
        });
      }
    });

    if (hovered || draggedTarget) {
      context.setLineDash([5, 5]);
      context.lineWidth = 1;

      [...texts].reverse().forEach((t, i) => {
        const x = t.x || 0;
        const y = t.y || 0;
        const width = t.width || MIN_WIDTH;
        const height = t.height || MIN_HEIGHT;

        // If box is hovered
        if ((texts.length - 1) - i == lastHover) {
          context.fillStyle = 'rgba(255, 255, 255, 0.2)';
          context.fillRect(x, y, width, height);
        }

        // Main box
        context.setLineDash([5, 5]);
        context.strokeStyle = 'white';
        context.strokeRect(x, y, width, height);

        context.setLineDash([0, 5, 5, 0]);
        context.strokeStyle = 'black';
        context.strokeRect(x, y, width, height);

        // Edges / Corners
        context.setLineDash([]);
        context.strokeStyle = 'black';
        context.fillStyle = 'rgba(255, 255, 255, 0.5)';

        // Top left corner
        context.fillRect(x-anchorSize, y-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect(x-anchorSize, y-anchorSize, anchorSize*2, anchorSize*2);

        // Top edge
        context.fillRect((x + (width/2))-anchorSize, y-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect((x + (width/2))-anchorSize, y-anchorSize, anchorSize*2, anchorSize*2);

        // Top right corner
        context.fillRect((x+width)-anchorSize, y-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect((x+width)-anchorSize, y-anchorSize, anchorSize*2, anchorSize*2);

        // Right edge
        context.fillRect((x+width)-anchorSize, (y+(height/2))-anchorSize, anchorSize*2, anchorSize*2)
        context.strokeRect((x+width)-anchorSize, (y+(height/2))-anchorSize, anchorSize*2, anchorSize*2)

        // Bottom right corner
        context.fillRect((x+width)-anchorSize, (y+height)-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect((x+width)-anchorSize, (y+height)-anchorSize, anchorSize*2, anchorSize*2);

        // Bottom edge
        context.fillRect((x + (width/2))-anchorSize, (y+height)-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect((x + (width/2))-anchorSize, (y+height)-anchorSize, anchorSize*2, anchorSize*2);

        // Bottom left corner
        context.fillRect(x-anchorSize, (y+height)-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect(x-anchorSize, (y+height)-anchorSize, anchorSize*2, anchorSize*2);

        // Left edge
        context.fillRect(x-anchorSize, (y+(height/2))-anchorSize, anchorSize*2, anchorSize*2);
        context.strokeRect(x-anchorSize, (y+(height/2))-anchorSize, anchorSize*2, anchorSize*2);
      });
    }
  }) as RenderFunc;

  let baked: Payload;
  const bake = async () => {
    console.log(JSON.stringify(texts.map(TemplateText.fromPartial)));

    const finalTexts = texts.map((t, i) => {
      const template = templateTexts[i];

      return Object.fromEntries(Object.entries(t).filter(([k, v]) => {
        if (k === 'template_text') {
          return false;
        }
        return v !== (template as any)[k];
      })) as Text;
    });

    const payload = Payload.fromPartial({
      version: 0,
      template: data.id,
      text: finalTexts
    });

    console.log(payload);

    baked = payload;

    const encoded = Payload.encode(payload).finish();
    const rawBase64 = btoa(String.fromCharCode(...encoded))
      .replace(/-/g, '+')
      .replace(/_/g, '/');

    console.log(Payload.decode(new TextEncoder().encode(atob(rawBase64))));

    const readableStream = new ReadableStream({
      start(controller) {
        controller.enqueue(encoded);
      },
      pull(controller) {
        controller.close();
      },
      cancel() {
        // Unused
      }
    });

    // eslint-disable-next-line no-undef
    const compressedReadableStream = readableStream.pipeThrough(new CompressionStream("gzip"));

    const buffer: Uint8Array[] = [];
    const out = new WritableStream({
      write: (chunk) => {
        buffer.push(chunk);
      }
    });

    compressedReadableStream.pipeTo(out).then(() => {
      let length = 0;
      buffer.forEach(item => length += item.length);

      let mergedArray = new Uint8Array(length);
      let offset = 0;
      buffer.forEach(item => {
        mergedArray.set(item, offset);
        offset += item.length;
      });

      let finalBase64 = btoa(String.fromCharCode(...mergedArray));
      if (finalBase64.length > rawBase64.length) {
        console.log('compression was too big');
        finalBase64 = rawBase64;
      }

      const padded = finalBase64.indexOf("=");
      if (padded >= 0) {
        finalBase64 = finalBase64.substring(0, padded);
      }

      copy(`${PUBLIC_BACKEND_API_BASE}/img/${finalBase64}.webp`);
      toastStore.trigger({
        message: "Copied to clipboard!"
      });
    });
  };

  const addText = () => {
    if (texts.length >= 20) {
      return;
    }

    const tempIndex = templateTexts.findIndex((_, i) => !texts.find(c => c.template_text == i));

    if (tempIndex >= 0) {
      texts = [...texts, {
        ...TemplateText.fromPartial(templateTexts[tempIndex]),
        template_text: tempIndex
      }];

      texts.sort((a, b) => (a.template_text === undefined ? Number.MAX_SAFE_INTEGER : a.template_text) - (b.template_text === undefined ? Number.MAX_SAFE_INTEGER : b.template_text));
    } else {
      texts = [...texts, Text.fromPartial({
        x: 0,
        y: 0,
        width: MIN_WIDTH,
        height: MIN_HEIGHT,
        fill_color: "#000000",
        stroke_color: "#000000"
      })];
    }
  };

  const removeText = (i: number) => {
    texts = [...texts.slice(0, i), ...texts.slice(i + 1)];
  };

  enum AnchorPosition {
    TopLeft,
    Top,
    TopRight,
    Right,
    BottomRight,
    Bottom,
    BottomLeft,
    Left
  }

  const calculateTarget = (e: MouseEvent): [number, AnchorPosition | undefined] | undefined => {
    const mouseX = e.offsetX;
    const mouseY = e.offsetY;

    for (let i = 0; i < texts.length; i++) {
      const t = texts[i];
      const x = t.x || 0;
      const y = t.y || 0;
      const width = t.width || MIN_WIDTH;
      const height = t.height || MIN_HEIGHT;

      // Check each corner/edge

      // Top left
      if (mouseX > (x - anchorSize) && mouseX < (x + anchorSize) && mouseY > (y - anchorSize) && mouseY < (y + anchorSize)) {
        return [i, AnchorPosition.TopLeft];
      }

      // Top
      if (mouseX > ((x + (width/2)) - anchorSize) && mouseX < ((x + (width/2)) + anchorSize) && mouseY > (y - anchorSize) && mouseY < (y + anchorSize)) {
        return [i, AnchorPosition.Top];
      }

      // Top Right
      if (mouseX > ((x+width) - anchorSize) && mouseX < ((x+width) + anchorSize) && mouseY > (y - anchorSize) && mouseY < (y + anchorSize)) {
        return [i, AnchorPosition.TopRight];
      }

      // Right
      if (mouseX > ((x+width) - anchorSize) && mouseX < ((x+width) + anchorSize) && mouseY > ((y+(height/2)) - anchorSize) && mouseY < ((y+(height/2)) + anchorSize)) {
        return [i, AnchorPosition.Right];
      }

      // Bottom Right
      if (mouseX > ((x+width) - anchorSize) && mouseX < ((x+width) + anchorSize) && mouseY > ((y+height) - anchorSize) && mouseY < ((y+height) + anchorSize)) {
        return [i, AnchorPosition.BottomRight];
      }

      // Bottom
      if (mouseX > ((x+(width/2)) - anchorSize) && mouseX < ((x+(width/2)) + anchorSize) && mouseY > ((y+height) - anchorSize) && mouseY < ((y+height) + anchorSize)) {
        return [i, AnchorPosition.Bottom];
      }

      // Bottom left
      if (mouseX > (x - anchorSize) && mouseX < (x + anchorSize) && mouseY > ((y+height) - anchorSize) && mouseY < ((y+height) + anchorSize)) {
        return [i, AnchorPosition.BottomLeft];
      }

      // Left
      if (mouseX > (x - anchorSize) && mouseX < (x + anchorSize) && mouseY > ((y+(height/2)) - anchorSize) && mouseY < ((y+(height/2)) + anchorSize)) {
        return [i, AnchorPosition.Left];
      }

      // Then check if we are inside the box.
      if ((mouseX > x && mouseX < x + width) && (mouseY > y && mouseY < y + height)) {
        return [i, undefined];
      }
    }
  }

  const anchorToCursor: Record<AnchorPosition, string> = {
    [AnchorPosition.TopLeft]: 'nwse-resize',
    [AnchorPosition.Top]: 'ns-resize',
    [AnchorPosition.TopRight]: 'nesw-resize',
    [AnchorPosition.Right]: 'ew-resize',
    [AnchorPosition.BottomRight]: 'nwse-resize',
    [AnchorPosition.Bottom]: 'ns-resize',
    [AnchorPosition.BottomLeft]: 'nesw-resize',
    [AnchorPosition.Left]: 'ew-resize',
  };


  let draggedTarget: undefined | [number, AnchorPosition | (undefined | null)];
  let dragRelative: [number, number];
  let dragSize: [number, number];
  let lastHover: number | undefined;
  const onMove = (e: MouseEvent) => {
    if (!draggedTarget) {
      const target = calculateTarget(e);
      if (target) {
        lastHover = target[0];
        if (target[1] !== undefined) {
          cursor = anchorToCursor[target[1]];
        } else {
          cursor = 'move';
        }
      } else {
        lastHover = undefined;
        cursor = 'initial';
      }
    } else {
      const h = texts[draggedTarget[0]].height || MIN_HEIGHT;
      const w = texts[draggedTarget[0]].width || MIN_HEIGHT;

      if (draggedTarget[1] !== undefined && draggedTarget[1] !== null) {
        const y = texts[draggedTarget[0]].y || 0;
        const x = texts[draggedTarget[0]].x || 0;

        switch (draggedTarget[1]) {
          case AnchorPosition.TopLeft:
          case AnchorPosition.Top:
          case AnchorPosition.TopRight:
            texts[draggedTarget[0]].height = Math.max(MIN_HEIGHT, dragSize[1] - e.offsetY);
            if (h > MIN_HEIGHT || (texts[draggedTarget[0]].height || MIN_HEIGHT) > MIN_HEIGHT) {
              texts[draggedTarget[0]].y = y - ((texts[draggedTarget[0]].height || MIN_HEIGHT) - h);
            }
            break;
        }

        switch (draggedTarget[1]) {
          case AnchorPosition.TopRight:
          case AnchorPosition.Right:
          case AnchorPosition.BottomRight:
            texts[draggedTarget[0]].width = Math.max(MIN_WIDTH, e.offsetX - x);
            break;
        }

        switch (draggedTarget[1]) {
          case AnchorPosition.BottomRight:
          case AnchorPosition.Bottom:
          case AnchorPosition.BottomLeft:
            texts[draggedTarget[0]].height = Math.max(MIN_HEIGHT, e.offsetY - y);
            break;
        }

        switch (draggedTarget[1]) {
          case AnchorPosition.BottomLeft:
          case AnchorPosition.Left:
          case AnchorPosition.TopLeft:
            texts[draggedTarget[0]].width = Math.max(MIN_WIDTH, dragSize[0] - e.offsetX);
            if (w > MIN_WIDTH || (texts[draggedTarget[0]].width || MIN_WIDTH) > MIN_WIDTH) {
              texts[draggedTarget[0]].x = x - ((texts[draggedTarget[0]].width || MIN_WIDTH) - w);
            }
            break;
        }

        texts[draggedTarget[0]].x = Math.max(0, Math.min(texts[draggedTarget[0]].x || 0, width - (texts[draggedTarget[0]].width || MIN_WIDTH)));
        texts[draggedTarget[0]].y = Math.max(0, Math.min(texts[draggedTarget[0]].y || 0, height - (texts[draggedTarget[0]].height || MIN_HEIGHT)));
      } else {
        texts[draggedTarget[0]].x = Math.min(Math.max(0, e.offsetX - dragRelative[0]), width - w);
        texts[draggedTarget[0]].y = Math.min(Math.max(0, e.offsetY - dragRelative[1]), height - h);
      }
    }
  }

  const onDown = (e: MouseEvent) => {
    draggedTarget = calculateTarget(e);
    if (draggedTarget) {
      const x = texts[draggedTarget[0]].x || 0;
      const y = texts[draggedTarget[0]].y || 0;
      dragRelative = [e.offsetX - x, e.offsetY - y];
      dragSize = [x + (texts[draggedTarget[0]].width || MIN_WIDTH), y + (texts[draggedTarget[0]].height || MIN_HEIGHT)]
    }
  }

  const onUp = () => {
    draggedTarget = undefined;
  }
</script>

<svelte:window on:mouseup={onUp}></svelte:window>

{#if loaded}
  <div class="flex flex-row justify-around gap-4 w-full max-w-screen-xl">
    <div class="w-fit h-fit select-none"
         on:mouseover={() => hovered = true}
         on:mouseout={() => hovered = false}
         on:focus={() => hovered = true}
         on:blur={() => hovered = false}
         on:mousemove={onMove}
         on:mousedown={onDown}
         style="cursor: {cursor}">
      <Canvas {width} {height}>
        <Layer {render} />
      </Canvas>
    </div>
    <div class="w-full h-full flex flex-col gap-3">
      <div>
        <button type="button" class="btn variant-ghost-success w-full" disabled={texts.length >= 20} on:click={() => addText()}>+</button>
      </div>
      {#each texts as text, i}
        <div class="flex flex-col">
          <div class="flex flex-row gap-2 mb-2">
            <textarea class="input p-2 w-full" bind:value={text.text}></textarea>
            <div class="flex flex-col gap-2 w-fit">
              <div class="flex flex-row gap-2 h-1/2">
                <SlideToggle name="slide" bind:checked={text.unfilled} />
                <input class="input" type="color" bind:value={text.fill_color}>
              </div>
              <div class="flex flex-row gap-2 h-1/2">
                <input type="number" min="1" max="72" class="input p-2" bind:value={text.stroke} />
                <input class="input" type="color" bind:value={text.stroke_color}>
              </div>
            </div>
          </div>
          <div class="flex flex-row gap-2">
            <select class="select" bind:value={text.font}>
              {#each FONTS as font}
                <option value={font}>{font}</option>
              {/each}
            </select>

            <input type="number" min="1" max="72" class="input p-2" bind:value={text.size} />

            <RadioGroup>
              <RadioItem bind:group={text.horizontal_align} name="justify" value={1}><Fa icon={faAlignLeft} class="text-2xl" /></RadioItem>
              <RadioItem bind:group={text.horizontal_align} name="justify" value={0}><Fa icon={faAlignCenter} class="text-2xl" /></RadioItem>
              <RadioItem bind:group={text.horizontal_align} name="justify" value={2}><Fa icon={faAlignRight} class="text-2xl" /></RadioItem>
            </RadioGroup>

            <RadioGroup>
              <RadioItem bind:group={text.vertical_align} name="justify" value={1}><Fa icon={faCircleUp} class="text-2xl" /></RadioItem>
              <RadioItem bind:group={text.vertical_align} name="justify" value={0}><Fa icon={faCircleDot} class="text-2xl" /></RadioItem>
              <RadioItem bind:group={text.vertical_align} name="justify" value={2}><Fa icon={faCircleDown} class="text-2xl" /></RadioItem>
            </RadioGroup>

            <button type="button" class="btn variant-ghost-error" on:click={() => removeText(i)}>X</button>
          </div>
          {#if i < texts.length - 1}
            <hr class="mt-3" />
          {/if}
        </div>
      {/each}
    </div>
  </div>
  <div class="mt-10 w-full max-w-screen-xl">
    <button type="button" class="btn variant-ghost-success w-full" on:click={() => bake()}>Bake</button>
  </div>
  <div class="mt-10 w-full max-w-screen-xl">
    <pre>{JSON.stringify(texts, null, 4)}</pre>
  </div>
{:else}
  Loading...
{/if}