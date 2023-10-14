import json

def parse_sarif(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)

    # Extract the SARIF version for reference (optional)
    sarif_version = data.get('version', '')

    # Extract the results from the SARIF file (assuming a single run)
    runs = data.get('runs', [])
    if not runs:
        return []  # No results found

    results = runs[0].get('results', [])

    # Initialize an empty list to store the parsed results
    parsed_results = []

    # Loop through each result in the results list
    for result in results:
        # Extract the necessary information from each result
        rule_id = result.get('ruleId', '')
        message = result.get('message', {}).get('text', '')
        
        # Extract location information if available
        locations = result.get('locations', [])
        if locations:
            location = locations[0].get('physicalLocation', {}).get('artifactLocation', {}).get('uri', '')
        else:
            location = ''

        # Append the extracted information to the parsed_results list
        parsed_results.append({
            'ruleId': rule_id,
            'message': message,
            'location': location
        })

    return {
        'sarif_version': sarif_version,
        'parsed_results': parsed_results
    }

# Example usage:
file_path = 'your_sarif_file.json'
parsed_data = parse_sarif(file_path)
print(parsed_data)
