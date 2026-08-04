document.addEventListener('DOMContentLoaded', () => {
  console.log('ready!');
  const form = document.querySelector('#convert');
  form.addEventListener('submit', handleForm);
});

async function handleForm(e) {
  e.preventDefault();
  const data = new FormData(e.currentTarget);
  try {
    const response = await fetch('/api/convert', {
      method: 'POST',
      headers: {
        // 'Content-Type': 'multipart/form-data'
      },
      body: data
    });
      if (!response.ok) {
        throw new Error(`${response.statusText}`);
      }

      const result = await response.json();

      const container = document.querySelector('#output');
      container.innerHTML = result;
  } catch (e) {
    console.error(e);
  }
}