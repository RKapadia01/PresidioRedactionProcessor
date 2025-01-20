import sys
import json
from presidio_anonymizer import AnonymizerEngine
from presidio_analyzer  import RecognizerResult

def main():
    input_data = sys.stdin.read()

    try:
        anonymizer_params = json.loads(input_data)

        # Parse the RecognizerResult objects from the JSON input
        analyzer_results = []
        for result in anonymizer_params["analyzer_results"]:
            analyzer_results.append(RecognizerResult.from_json(result))

        anonymizer = AnonymizerEngine()
        results = anonymizer.anonymize(text=anonymizer_params["text"], analyzer_results=analyzer_results)

        print(results.text)
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON input. {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()