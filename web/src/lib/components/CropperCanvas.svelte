<script lang="ts">
        import { onMount, createEventDispatcher } from 'svelte';
        import Cropper from 'cropperjs';
        import 'cropperjs/dist/cropper.css';

        export let src: string | null = null;
        let imageEl: HTMLImageElement;
        let cropper: Cropper | null = null;
        const dispatch = createEventDispatcher();

        onMount(() => {
          if (imageEl && src) {
            cropper = new Cropper(imageEl, {
              aspectRatio: 1,
              viewMode: 1,
              autoCropArea: 1,
              responsive: true,
              background: false,
              minContainerWidth: 300,
              minContainerHeight: 300,
              minCanvasWidth: 300,
              minCanvasHeight: 300,
              ready() {
                dispatch('ready');
              }
            });
          }
          return () => cropper?.destroy();
        });

        export async function toBlob(options?: { type?: string; quality?: number }) {
          return new Promise<Blob>((resolve, reject) => {
            if (!cropper) return reject('Cropper not initialized');
            cropper.getCroppedCanvas().toBlob(
              (blob) => {
                if (blob) resolve(blob);
                else reject('Failed to crop');
              },
              options?.type || 'image/png',
              options?.quality || 1
            );
          });
        }
      </script>

      {#if src}
        <div class="cropper-container">
          <img bind:this={imageEl} src={src} alt="Cropper" />
        </div>
      {:else}
        <div class="text-gray-400">No image selected</div>
      {/if}

      <style>
        .cropper-container {
          max-width: 350px;
          max-height: 350px;
          margin: 0 auto;
        }
        .cropper-container img {
          max-width: 100%;
          max-height: 350px;
          display: block;
        }
      </style>