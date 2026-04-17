import { state } from '../../core/state.js';
import { renderLoading, renderError } from '../../shared/render.js';

// renderUpload displays the form for uploading images.
function renderUpload() {
  return `
    <div id="filters">
      <div class="filters-header">
        <h4>Upload Media</h4>
        <button id="toggle-upload-form">Toggle</button>
      </div>

      <div class="controls ${state.showcase.ui.formOpen ? '' : 'hidden'}">

        <form id="upload-form">

          <div class="text-search">
            <div class="text-group">
              <label for="media-file">Image</label>
              <input id="media-file" type="file" accept="image/*" required />
            </div>
          </div>

          <div class="select-group">
            <label for="media-filter">Filter</label>
            <select id="media-filter">
              <option value="grayscale">Grayscale</option>
            </select>
          </div>

          <button id="apply-upload" type="submit">Upload</button>
        </form>

      </div>
    </div>
  `;
}

// renderImagesGrid displays showcase images.
function renderImagesGrid() {
  if (state.showcase.data.length === 0) {
    return `<p class="empty-state">No media uploaded yet.</p>`;
  }

  return `
    <section class="grid">
      ${state.showcase.data
        .map(
          (job) => `
            <div class="media-card" data-job-id="${job.id}">
              <div class="media-preview">
                <img src="${job.input_url}" />
              </div>

              <div class="media-meta">
                <p></p>
                <p></p>
              </div>
            </div>
          `,
        )
        .join('')}
    </section>
  `;
}

// renderModal displays a modal for a specific image.
function renderModal() {
  const job = state.showcase.ui.selectedJob;
  if (!state.showcase.ui.modalOpen || !job) return '';

  const processedContent = job.output_url
    ? `<img src="${job.output_url}" />`
    : renderLoading('Processing image...');

  return `
    <div class="modal-overlay">
      <div class="modal">
        <button id="modal-close">Close</button>

        <h3>Media Details</h3>

        <div class="modal-body">
          <div>
            <h4>Original</h4>
            <img src="${job.input_url}" />
          </div>

          <div>
            <h4>Processed</h4>
            ${processedContent}
          </div>
        </div>

        <div class="modal-meta">
        </div>
      </div>
    </div>
  `;
}

// render compose the html for showcase.html.
export function render() {
  const app = document.querySelector('#app');

  let content = '';

  if (state.showcase.loading) {
    content = renderLoading('Loading media...');
  } else if (state.showcase.error) {
    content = renderError(state.showcase.error);
  } else {
    content = `
      <h2>Hotel Showcase</h2>
      ${renderUpload()}
      ${renderImagesGrid()}
      ${renderModal()}
    `;
  }

  app.innerHTML = content;
}
