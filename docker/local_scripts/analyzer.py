import sys
import json
from presidio_analyzer import AnalyzerEngine
from presidio_anonymizer import AnonymizerEngine

def main():
    input_data = sys.stdin.read()

    try:
        analyzer_params = json.loads(input_data)
        analyzer = AnalyzerEngine()
        results = analyzer.analyze(**analyzer_params)

        # Create the response array
        response = []
        for res in results:
            response.append(res.to_dict())

        print(json.dumps(response))
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON input. {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()