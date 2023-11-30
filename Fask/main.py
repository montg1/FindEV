from flask import Flask, render_template, request

app = Flask(__name__)

@app.route('/calculate_route', methods=['POST'])
def calculate_route():
    series = request.form['series']
    current_location = request.form['current_location']
    destination = request.form['destination']

    # Perform calculations or any other processing here

    return render_template('result.html', series=series, current=current_location, destination=destination)

@app.route('/')
def index():
    return render_template('index.html')

if __name__ == '__main__':
    app.run(debug=True)

