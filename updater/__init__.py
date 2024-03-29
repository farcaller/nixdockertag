import subprocess
import json
import hashlib
import os
import glob
import urllib

import typer
from .dxf.dxf import DXF
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

  def auth(d, response):
    d.authenticate(response=response)
  d = DXF(host, repo, auth)
  _, dcd = d._get_alias(
    alias=info['followTag'],
    manifest=None,
    verify=True,
    sizes=False,
    get_digest=False,
    get_dcd=True,
    get_manifest=False,
    platform=None,
    ml=True,
  )
  hash = dcd.split(':')[1]

  if hash == info['hash']:
    return
  
  print(f'updating {name} to {hash}')
  
  image_path = os.path.join('images', f'{name}.nix')

  with open(image_path, 'w') as f:
    f.write(chevron.render(
      IMAGE_TEMPLATE,
      data=dict(
        host=host,
        repo=repo,
        follow=info['followTag'],
        hash=hash)))
  
  subprocess.run(['git', 'add', image_path], check=True)
  
  if commit:
    subprocess.run(['git', 'commit', '-m', f'{info["image"]}: update to {hash}'], check=True)

@app.command()
def update_all(commit: bool = typer.Option(False)):
  images = glob.glob('images/*.nix')
  for name in images:
    name = os.path.splitext(os.path.basename(name))[0]
    print(f'checking {name}')
    try:
      update(name, commit)
    except RuntimeError as e:
      print(f'failed: {e}')

if __name__ == "__main__":
  app()
