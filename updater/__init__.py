import subprocess
import json
import hashlib
import os

import typer
from dxf import DXF
import requests
import chevron

app = typer.Typer()

IMAGE_TEMPLATE = '''{
  image = "{{ host }}/{{ repo }}";
  followTag = "{{ follow }}";
  hash = "{{ hash }}";
}
'''

def get_image_info(name: str):
  return subprocess.run([
    'nix', 'eval', '--expr', f'builtins.toJSON (import ./images/{name}.nix)', '--impure', '--raw'
  ], capture_output=True, text=True)

@app.command()
def update(name: str, commit: bool = typer.Option(False)):
  info = json.loads(get_image_info(name).stdout)
  host, repo = info['image'].split('/', 1)

  def auth(dxf, response):
    token = requests.get(f'https://ghcr.io/token?service=ghcr.io&scope=repository:{repo}:pull&client_id=updater').json()
    dxf.authenticate(response=response, authorization=f'Bearer{token["token"]}')

  d = DXF(host, repo, auth)
  mf = d.get_manifest(info['followTag'])
  hash = hashlib.sha256(mf.encode('utf8')).hexdigest();
  
  image_path = os.path.join('images', f'{name}.nix')

  with open(image_path, 'w') as f:
    f.write(chevron.render(
      IMAGE_TEMPLATE,
      data=dict(
        host=host,
        repo=repo,
        follow=info['followTag'],
        hash=hash)))
  
  if commit:
    subprocess.run(['git', 'add', image_path], check=True)
    subprocess.run(['git', 'commit', '-m', f'{info["image"]}: update to {hash}'], check=True)

if __name__ == "__main__":
  app()
