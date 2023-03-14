# -*- coding: utf-8 -*-
import os
import numpy as np
import pickle as pkl
import tensorflow as tf
from tensorflow.contrib import learn

import data_loader

# Show warnings and errors only
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3' # 屏蔽所有级别的日志
tf.compat.v1.logging.set_verbosity(tf.compat.v1.logging.ERROR) # 设置日志等级为ERROR

# File paths
tf.flags.DEFINE_string('test_data_file', None, '''Test data file path''')

# Test batch size
tf.flags.DEFINE_integer('batch_size', 1, 'Test batch size')
FLAGS = tf.flags.FLAGS


def traffic_detect(urls):
    # Restore parameters
    with open(os.path.join('params/params.pkl'), 'rb') as f:
        params = pkl.load(f, encoding='bytes')

    # Restore vocabulary processor
    vocab_processor = learn.preprocessing.VocabularyProcessor.restore(os.path.join('params/vocab'))

    # Load test data
    data, labels, lengths, _ = data_loader.load_data(urls=urls,
                                                     sw_path=params['stop_word_file'],
                                                     min_frequency=params['min_frequency'],
                                                     max_length=params['max_length'],
                                                     language=params['language'],
                                                     vocab_processor=vocab_processor,
                                                     shuffle=False)

    # Restore graph
    graph = tf.Graph()
    with tf.Session(graph=tf.Graph()) as sess:
        sess = tf.Session()

        tf.saved_model.loader.load(sess, ['serve'], 'model')
        graph = tf.get_default_graph()

        # Get tensors
        input_x = graph.get_tensor_by_name('input_x:0')
        input_y = graph.get_tensor_by_name('input_y:0')
        keep_prob = graph.get_tensor_by_name('keep_prob:0')
        predictions = graph.get_tensor_by_name('softmax/predictions:0')
        accuracy = graph.get_tensor_by_name('accuracy/accuracy:0')

        # Generate batches
        batches = data_loader.batch_iter(data, labels, lengths, FLAGS.batch_size, 1)

        all_predictions = []
        # num_batches = int(len(data) / FLAGS.batch_size)
        # sum_accuracy = 0

        # Test
        for batch in batches:
            x_test, y_test, x_lengths = batch
            if params['clf'] == 'cnn':
                feed_dict = {input_x: x_test, input_y: y_test, keep_prob: 1.0}
                batch_predictions, batch_accuracy = sess.run([predictions, accuracy], feed_dict)
            else:
                batch_size = graph.get_tensor_by_name('batch_size:0')
                sequence_length = graph.get_tensor_by_name('sequence_length:0')
                feed_dict = {input_x: x_test, input_y: y_test, batch_size: FLAGS.batch_size, sequence_length: x_lengths,
                             keep_prob: 1.0}

                batch_predictions, batch_accuracy = sess.run([predictions, accuracy], feed_dict)
            # sum_accuracy += batch_accuracy
            all_predictions = np.concatenate([all_predictions, batch_predictions])

        # final_accuracy = sum_accuracy / num_batches

    # Print test accuracy
    # print('Test accuracy: {}'.format(final_accuracy))

    return all_predictions

if __name__ == "__main__":
    urls = ["get /questions/30562", "get /questions/30562", "get /questions/30562", "get /questions/30562"]
    print(traffic_detect(urls))
