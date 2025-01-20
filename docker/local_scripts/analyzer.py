import sys
import json
from presidio_analyzer import AnalyzerEngine

def main():
    input_data = sys.stdin.read()

    try:
        analyzer_inputs = json.loads(input_data)
        analyzer = AnalyzerEngine()
        analyzer_results = analyzer.analyze(**analyzer_inputs)

        # Create the response array
        response = []
        for analyzer_result in analyzer_results:
            response.append(analyzer_result.to_dict())

        print(json.dumps(response))
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON input. {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()