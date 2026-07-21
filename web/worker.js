let ready = false;
const queued = [];

const go = new Go();

WebAssembly.instantiateStreaming(fetch('paa.wasm'), go.importObject)
  .then((result) => {
    go.run(result.instance).catch((err) => {
      postMessage({ type: 'fatal', error: String(err) });
    });

    ready = true;
    postMessage({ type: 'ready' });
    for (const msg of queued.splice(0)) {
      handleConvert(msg);
    }
  })
  .catch((err) => {
    postMessage({ type: 'fatal', error: 'failed to load paa.wasm: ' + err });
  });

onmessage = (event) => {
  if (!ready) {
    queued.push(event.data);
    return;
  }
  handleConvert(event.data);
};

function handleConvert(data) {
  let result;
  try {
      result = paaToPng(new Uint8Array(data));
  } catch (err) {
    result = { ok: false, error: String(err) };
  }

  if (!result.ok) {
    postMessage({ ok: false, error: result.error });
    return;
  }

  postMessage(
    { ok: true, png: result.png.buffer },
    [result.png.buffer]
  );
}
