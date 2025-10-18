import csv

component_names = {
    'SERVER': 'server',
    'FRONTEND': 'frontend',
    'DB': 'database',
    'EXPORTER': 'exporter',
    'SIMULATOR': 'simulator',
    'CI': 'CI/CD',
    'BRUNO': 'Bruno tests',
    'META': 'project files'
}

def get_tag(path):
    if path.startswith('server/'):
        return 'SERVER'
    elif path.startswith('frontend/'):
        return 'FRONTEND'
    elif path.startswith('migrations/'):
        return 'DB'
    elif path.startswith('exporter/'):
        return 'EXPORTER'
    elif path.startswith('simulator/'):
        return 'SIMULATOR'
    elif path.startswith('.github/'):
        return 'CI'
    elif path.startswith('WanBingo Bruno/'):
        return 'BRUNO'
    else:
        return 'META'

with open('/tmp/commit_log.txt', 'r') as f:
    lines = f.readlines()

commits = []
current_commit = None
files = []
for line in lines:
    line = line.strip()
    if not line:
        if current_commit:
            commits.append((current_commit, files))
            files = []
            current_commit = None
        continue
    if '|' in line:
        if current_commit:
            commits.append((current_commit, files))
            files = []
        hash_msg = line.split('|', 1)
        current_commit = (hash_msg[0], hash_msg[1])
    else:
        files.append(line)

if current_commit:
    commits.append((current_commit, files))

with open('new-names.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile)
    writer.writerow(['commit_hash', 'original_message', 'new_name'])
    for commit, files in commits:
        hash_, msg = commit
        tags = set()
        for file in files:
            tag = get_tag(file)
            tags.add(tag)
        if not tags:
            continue
        components = [component_names[tag] for tag in sorted(tags)]
        if len(components) == 1:
            summary = f"Update {components[0]}"
        else:
            summary = f"Update {', '.join(components[:-1])} and {components[-1]}"
        tags_str = ','.join(sorted(tags))
        new_name = f"[{tags_str}] {summary}"
        writer.writerow([hash_, msg, new_name])