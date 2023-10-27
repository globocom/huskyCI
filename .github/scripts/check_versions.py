import yaml
import requests


def get_latest_docker_image_version(repo_name):
    url = f"https://registry.hub.docker.com/v2/repositories/huskyci/{repo_name}/tags"
    response = requests.get(url)
    data = response.json()

    # Assuming the first result is the latest
    latest_version = data['results'][0]['name']
    return latest_version


def main():
    with open('config.yaml', 'r') as f:
        config = yaml.safe_load(f)

    for tool, tool_info in config.items():
        current_version = tool_info['imageTag']
        latest_version = get_latest_docker_image_version(
            tool_info['image'].split('/')[-1])

        if current_version != latest_version:
            print(
                f"[WARNING] {tool} is outdated. Current: {current_version}, Latest: {latest_version}")
        else:
            print(f"{tool} is up-to-date with version {current_version}.")


if __name__ == "__main__":
    main()
